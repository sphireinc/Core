package core

import (
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

const defaultLimiterWindow = 5 * time.Second
const defaultLimiterMaxRequests = 10

func clientIdentity(ctx *Context) string {
	if ctx == nil {
		return ""
	}

	if forwarded := string(ctx.Request.Header.Peek("X-Forwarded-For")); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if realIP := string(ctx.Request.Header.Peek("X-Real-IP")); realIP != "" {
		return strings.TrimSpace(realIP)
	}

	if ip := ctx.RequestCtx.RemoteIP(); ip != nil {
		return ip.String()
	}

	return fmt.Sprintf("conn:%d", ctx.ConnID())
}

func limiterKey(identity string) uint64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(identity))
	return hasher.Sum64()
}

func limiterWindow(cfg *Config) time.Duration {
	if cfg == nil {
		return defaultLimiterWindow
	}
	return cfg.Middleware.Limiter.Time.Or(defaultLimiterWindow)
}

func limiterMax(cfg *Config) uint64 {
	if cfg == nil || cfg.Middleware.Limiter.Max <= 0 {
		return defaultLimiterMaxRequests
	}
	return uint64(cfg.Middleware.Limiter.Max)
}

func writeJSONError(ctx *Context, status int, body string) error {
	ctx.SetStatusCode(status)
	ctx.SetContentType("application/json")
	ctx.SetBody([]byte(body))
	return nil
}

// Limiter creates a basic in-memory limiter
func Limiter(ctx *Context) error {
	cfg := configFromContext(ctx)
	if cfg == nil {
		return nil
	}

	store := cfg.memoryStore()
	if store == nil {
		return nil
	}

	identity := clientIdentity(ctx)
	key := limiterKey(identity)

	window := limiterWindow(cfg)
	maxRequests := limiterMax(cfg)

	var requestCount uint64 = 1
	if stored, ok := store.Get(key); ok && stored != nil {
		if count, castOK := stored.(uint64); castOK {
			requestCount = count + 1
		}
	}

	if requestCount > maxRequests {
		ctx.Response.Header.Set("Retry-After", fmt.Sprintf("%.0f", window.Seconds()))
		return writeJSONError(ctx, fasthttp.StatusTooManyRequests, `{"status":"error","error":"too many requests"}`)
	}

	store.Set(key, requestCount, time.Now().Add(window))
	return nil
}

// BasicAuth checks for basic authentication parameters.
//
// This implementation is intentionally conservative:
// it validates a syntactically correct HTTP Basic Authorization header,
// but does not invent a new credential configuration schema.
func BasicAuth(ctx *Context) error {
	authHeader := strings.TrimSpace(string(ctx.Request.Header.Peek("Authorization")))
	if authHeader == "" {
		ctx.Response.Header.Set("WWW-Authenticate", `Basic realm="Core"`)
		return writeJSONError(ctx, fasthttp.StatusUnauthorized, `{"status":"error","error":"missing authorization header"}`)
	}

	const prefix = "Basic "
	if !strings.HasPrefix(authHeader, prefix) {
		ctx.Response.Header.Set("WWW-Authenticate", `Basic realm="Core"`)
		return writeJSONError(ctx, fasthttp.StatusUnauthorized, `{"status":"error","error":"invalid authorization scheme"}`)
	}

	encoded := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
	if encoded == "" {
		ctx.Response.Header.Set("WWW-Authenticate", `Basic realm="Core"`)
		return writeJSONError(ctx, fasthttp.StatusUnauthorized, `{"status":"error","error":"missing basic credentials"}`)
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		ctx.Response.Header.Set("WWW-Authenticate", `Basic realm="Core"`)
		return writeJSONError(ctx, fasthttp.StatusUnauthorized, `{"status":"error","error":"invalid basic credentials"}`)
	}

	credentials := string(decoded)
	separator := strings.IndexByte(credentials, ':')
	if separator < 0 {
		ctx.Response.Header.Set("WWW-Authenticate", `Basic realm="Core"`)
		return writeJSONError(ctx, fasthttp.StatusUnauthorized, `{"status":"error","error":"malformed basic credentials"}`)
	}

	username := strings.TrimSpace(credentials[:separator])
	if username == "" {
		ctx.Response.Header.Set("WWW-Authenticate", `Basic realm="Core"`)
		return writeJSONError(ctx, fasthttp.StatusUnauthorized, `{"status":"error","error":"invalid basic credentials"}`)
	}

	return nil
}

// AdminOnly checks if a user is an admin.
//
// This trusts upstream identity propagation through one of:
//   - X-Admin: true|1|yes|on
//   - X-Role: admin
//   - X-Roles: admin[, ...]
func AdminOnly(ctx *Context) error {
	if isTruthyHeader(ctx, "X-Admin") {
		return nil
	}

	role := strings.TrimSpace(strings.ToLower(string(ctx.Request.Header.Peek("X-Role"))))
	if role == "admin" {
		return nil
	}

	rolesHeader := string(ctx.Request.Header.Peek("X-Roles"))
	if rolesHeader != "" {
		for _, rolePart := range strings.Split(rolesHeader, ",") {
			if strings.TrimSpace(strings.ToLower(rolePart)) == "admin" {
				return nil
			}
		}
	}

	return writeJSONError(ctx, fasthttp.StatusForbidden, `{"status":"error","error":"admin access required"}`)
}

func isTruthyHeader(ctx *Context, key string) bool {
	value := strings.TrimSpace(strings.ToLower(string(ctx.Request.Header.Peek(key))))
	switch value {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

// LogRequest provides a log middleware which logs each request.
func LogRequest(ctx *Context) error {
	cfg := configFromContext(ctx)
	if cfg != nil {
		logger := cfg.logger()
		if logger != nil {
			logger.Info(fmt.Sprintf("%s %s", ctx.RequestCtx.RemoteAddr(), ctx.Request.URI()))
		}
	}
	return nil
}

// BasicHeaders applies our general headers.
func BasicHeaders(ctx *Context) error {
	cfg := configFromContext(ctx)

	if len(ctx.Response.Header.ContentType()) == 0 {
		ctx.Response.Header.Set("Content-Type", "application/json")
	}
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Pragma", "no-cache")

	if cfg != nil && cfg.Middleware.BasicHeaders.ShowServer {
		ctx.Response.Header.Set("Server", "Core")
	}

	origin := string(ctx.Request.Header.Peek("Origin"))
	if origin != "" {
		requestHeaders := string(ctx.Request.Header.Peek("Access-Control-Request-Headers"))
		if requestHeaders == "" {
			requestHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
		}

		requestMethod := string(ctx.Request.Header.Peek("Access-Control-Request-Method"))
		if requestMethod == "" {
			requestMethod = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
		}

		ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", requestMethod)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", requestHeaders)
		ctx.Response.Header.Set("Vary", "Origin")
		ctx.Response.Header.Add("Vary", "Access-Control-Request-Method")
		ctx.Response.Header.Add("Vary", "Access-Control-Request-Headers")
	}

	return nil
}
