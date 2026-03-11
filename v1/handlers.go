package core

import (
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

func applyResponseHeaders(ctx *Context, cfg *Config) {
	if cfg == nil {
		return
	}

	if cfg.Server.ResponseConfig.SessionToken {
		if sessionToken := HeaderFromCtx(ctx, "Session-Token"); sessionToken != "" {
			ctx.Response.Header.Set("X-Session-Token", sessionToken)
		}
	}

	if cfg.Server.ResponseConfig.RequestID {
		requestID := HeaderFromCtx(ctx, "Request-Id")
		if requestID == "" {
			requestID = NextRequestID()
		}
		ctx.Response.Header.Set("X-Request-Id", requestID)
	}
}

func responseMetadataFromContext(ctx *Context, cfg *Config) Metadata {
	metadata := Metadata{
		RequestTime: time.Now().UTC().Format(time.RFC3339Nano),
	}

	if cfg == nil {
		return metadata
	}

	if cfg.Server.ResponseConfig.SessionToken {
		metadata.SessionToken = HeaderFromCtx(ctx, "Session-Token")
	}

	if cfg.Server.ResponseConfig.RequestID {
		metadata.RequestID = HeaderFromCtx(ctx, "Request-Id")
		if metadata.RequestID == "" {
			metadata.RequestID = NextRequestID()
		}
	}

	return metadata
}

// ResponseJson creates a body and returns a JSON response.
func ResponseJson(ctx *Context, output string, status int, errorText string) error {
	cfg := configFromContext(ctx)

	if output != "" {
		return HandleResponseJSON(ctx, []byte(output), status)
	}

	body := responseBytes(
		"error",
		nil,
		&APIError{
			Code:    "request_error",
			Message: errorText,
		},
		responseMetadataFromContext(ctx, cfg),
	)

	return HandleResponseJSON(ctx, body, status)
}

// HandleResponseJSON handles general responses via JSON.
func HandleResponseJSON(ctx *Context, body []byte, status int) error {
	ctx.SetContentType("application/json")

	cfg := configFromContext(ctx)
	applyResponseHeaders(ctx, cfg)

	ctx.SetStatusCode(fasthttp.StatusOK)
	if status != fasthttp.StatusOK {
		ctx.SetStatusCode(status)
	}

	ctx.SetBody(body)
	return nil
}

// NotFoundServer is the default 404 handler.
func NotFoundServer(ctx *Context) error {
	cfg := configFromContext(ctx)
	body := responseBytes(
		"error",
		nil,
		&APIError{
			Code:    "not_found",
			Message: "resource not found",
			Details: string(ctx.Request.RequestURI()),
		},
		responseMetadataFromContext(ctx, cfg),
	)
	return HandleResponseJSON(ctx, body, fasthttp.StatusNotFound)
}

func MethodNotAllowed(ctx *Context) error {
	cfg := configFromContext(ctx)
	body := responseBytes(
		"error",
		nil,
		&APIError{
			Code:    "method_not_allowed",
			Message: "method not allowed",
			Details: string(ctx.Request.RequestURI()),
		},
		responseMetadataFromContext(ctx, cfg),
	)
	return HandleResponseJSON(ctx, body, fasthttp.StatusMethodNotAllowed)
}

// Status simply returns a 200 OK.
func Status(ctx *Context) error {
	cfg := configFromContext(ctx)
	body := responseBytes(
		"ok",
		map[string]string{"message": "OK"},
		nil,
		responseMetadataFromContext(ctx, cfg),
	)
	return HandleResponseJSON(ctx, body, fasthttp.StatusOK)
}

// TeaPot is a 418 handler easter egg.
func TeaPot(ctx *Context) error {
	ctx.Response.Header.Set("X-Teapot", "Chai")

	cfg := configFromContext(ctx)
	body := responseBytes(
		"error",
		map[string]string{"message": "Are you a teapot?"},
		&APIError{
			Code:    "teapot",
			Message: "request refused by teapot",
		},
		responseMetadataFromContext(ctx, cfg),
	)
	return HandleResponseJSON(ctx, body, fasthttp.StatusTeapot)
}

// HeaderFromCtx returns the requested header from the given fasthttp context.
func HeaderFromCtx(ctx *Context, key string) string {
	return string(ctx.Request.Header.Peek(key))
}

// NextRequestID returns the next request id, which is a unix timestamp.
func NextRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
