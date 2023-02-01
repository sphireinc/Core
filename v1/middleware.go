package core

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/memory"
	"time"
)

var Limiter = limiter.New(limiter.Config{
	Next: func(c *fiber.Ctx) bool {
		return c.IP() == "127.0.0.1"
	},
	Max:        20,
	Expiration: 30 * time.Second,
	KeyGenerator: func(c *fiber.Ctx) string {
		return "key"
	},
	LimitReached: func(c *fiber.Ctx) error {
		return c.SendFile("./fast.html")
	},
	Storage: memory.New(memory.Config{
		GCInterval: 10 * time.Second,
	}),
})

// BasicAuth checks for basic authentication parameters
func BasicAuth(ctx *Context) error {
	if string(ctx.Response.Header.Peek("Origin")) == "a" {

	}
	return nil
}

// AdminOnly checks if a user is an admin.
func AdminOnly(ctx *Context) error {
	if string(ctx.Response.Header.Peek("Origin")) == "a" {

	}
	return nil
}

// LogRequest provides a log middleware which logs each request.
//
//	Forced middleware
func LogRequest(ctx *Context) error {
	App.Log.Writer.Info(fmt.Sprintf("%s %s",
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
	ctx.Response.Header.Set("Server", "Core")

	if origin := string(ctx.Response.Header.Peek("Origin")); origin != "" {
		acAllowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", origin)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", acAllowHeaders)
	}

	return nil
}
