package core

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Limiter creates a basic in-memory limiter
func Limiter(ctx *Context) error {
	val, ok := App.Persistence.Memory.Get(ctx.ConnID())

	newVal := uint64(1)
	if ok && val != nil {
		newVal = val.(uint64) + uint64(1)
	}

	if App.Middleware.Limiter.Time == 0 {
		App.Middleware.Limiter.Time = 5
	}
	newTime := time.Now().Add(time.Duration(App.Middleware.Limiter.Time))

	if App.Middleware.Limiter.Max == 0 {
		App.Middleware.Limiter.Max = 10
	}

	if newVal > uint64(App.Middleware.Limiter.Max) {
		return errors.New("capacity limit reached id: " + strconv.FormatUint(ctx.ConnID(), 10))
	}

	App.Persistence.Memory.Set(ctx.ConnID(), newVal, newTime)
	return nil
}

// BasicAuth checks for basic authentication parameters
func BasicAuth(ctx *Context) error {
	//if string(ctx.Response.Header.Peek("Origin")) == "a" {
	//
	//}
	return nil
}

// AdminOnly checks if a user is an admin.
func AdminOnly(ctx *Context) error {
	//if string(ctx.Response.Header.Peek("Origin")) == "a" {
	//
	//}
	return nil
}

// LogRequest provides a log middleware which logs each request.
//
//	Forced middleware
func LogRequest(ctx *Context) error {
	defer App.Log.Writer.Info(fmt.Sprintf("%s %s",
		ctx.RequestCtx.RemoteAddr(),
		ctx.Request.URI()))
	return nil
}

// BasicHeaders applies our general headers.
// Forced middleware
func BasicHeaders(ctx *Context) error {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.Header.Set("Accept", "*/*")
	ctx.Response.Header.Set("Accept-Encoding", "gzip, compress, deflate, br")
	ctx.Response.Header.Set("Accept-Language", "en-US,en;q=0.8")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Pragma", "no-cache")
	ctx.Response.Header.Set("Upgrade-Insecure-Requests", "1")

	if App.Middleware.BasicHeaders.ShowServer {
		ctx.Response.Header.Set("Server", "Core")
	}

	if origin := string(ctx.Response.Header.Peek("Origin")); origin != "" {
		acAllowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", acAllowHeaders)
	}

	return nil
}
