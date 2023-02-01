package core

import (
	"errors"
	mantisHttp "github.com/sphireinc/mantis/http"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

// HandleResponseJSON handles general responses via JSON.
func HandleResponseJSON(ctx *Context, body []byte, status int) error {
	ctx.SetContentType("application/json")

	// Set session token and request ID if available
	ctx.Response.Header.Set("X-Session-Token", HeaderFromCtx(ctx, "Session-Token"))
	ctx.Response.Header.Set("X-Request-Id", HeaderFromCtx(ctx, "Request-Id"))

	ctx.SetStatusCode(fasthttp.StatusOK)
	if status != 200 {
		ctx.SetStatusCode(status)
	}

	ctx.SetBody(body)
	return nil
}

// NotFoundServer is the default 404 handler
func NotFoundServer(ctx *Context) error {
	body := mantisHttp.Response{
		Body:  []byte("404 Not Found: " + string(ctx.Request.RequestURI())),
		Error: errors.New("404 Not Found"),
	}
	return HandleResponseJSON(ctx, body.Byte(), fasthttp.StatusNotFound)
}

func MethodNotAllowed(ctx *Context) error {
	body := mantisHttp.Response{
		Body:  []byte("405 Method Not Allowed" + string(ctx.Request.RequestURI())),
		Error: errors.New("405 Method Not Allowed"),
	}
	return HandleResponseJSON(ctx, body.Byte(), fasthttp.StatusOK)
}

// Status simply returns a 200 OK
func Status(ctx *Context) error {
	body := mantisHttp.Response{Body: []byte("OK")}
	return HandleResponseJSON(ctx, body.Byte(), fasthttp.StatusOK)
}

// TeaPot is a 418 handler easter egg
func TeaPot(ctx *Context) error {
	ctx.Response.Header.Set("X-Teapot", "Chai")
	body := mantisHttp.Response{Body: []byte("Are you a teapot?")}
	return HandleResponseJSON(ctx, body.Byte(), fasthttp.StatusTeapot)
}

// HeaderFromCtx returns the requested header from the given fasthttp context
func HeaderFromCtx(ctx *Context, key string) string {
	return string(ctx.Request.Header.Peek(key))
}

// NextRequestID returns the next request id, which is a unix timestamp.
func NextRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}
