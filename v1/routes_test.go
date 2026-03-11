package core

import (
	"testing"

	routing "github.com/qiangxue/fasthttp-routing"
)

func testHandler(ctx *Context) error { return nil }

func TestMethodsReturnsDistinctMethods(t *testing.T) {
	cfg := &Config{}
	methods := cfg.Methods()

	expected := []string{"GET", "PUT", "POST", "DELETE"}
	if len(methods) != len(expected) {
		t.Fatalf("expected %d methods, got %d", len(expected), len(methods))
	}

	for i, method := range expected {
		if methods[i] != method {
			t.Fatalf("expected %q at index %d, got %q", method, i, methods[i])
		}
	}
}

func TestAttachMiddlewareIncludesRouteSpecificMiddleware(t *testing.T) {
	cfg := &Config{}
	cfg.Middleware.LogRequest.Enabled = true
	cfg.Middleware.BasicHeaders.Enabled = true

	routeMW := []routing.Handler{
		func(ctx *Context) error { return nil },
	}

	handlers := attachMiddleware(cfg, routeMW, testHandler)

	// injectConfig + LogRequest + BasicHeaders + route middleware + final handler
	if len(handlers) != 5 {
		t.Fatalf("expected 5 handlers, got %d", len(handlers))
	}
}

func TestGroupPrefixesRoute(t *testing.T) {
	router := &Router{}
	group := router.Group("/api/v1")
	group.Get("/status", testHandler)

	if len(router.Routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(router.Routes))
	}
	if router.Routes[0].URI != "/api/v1/status" {
		t.Fatalf("expected URI /api/v1/status, got %q", router.Routes[0].URI)
	}
}

func TestVersionPrefixesRoute(t *testing.T) {
	router := &Router{}
	versioned := router.Version("v2")
	versioned.Post("widgets", testHandler)

	if len(router.Routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(router.Routes))
	}
	if router.Routes[0].URI != "/v2/widgets" {
		t.Fatalf("expected URI /v2/widgets, got %q", router.Routes[0].URI)
	}
	if router.Routes[0].Method != "POST" {
		t.Fatalf("expected method POST, got %q", router.Routes[0].Method)
	}
}
