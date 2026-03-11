package core

import (
	"fmt"
	"reflect"

	routing "github.com/qiangxue/fasthttp-routing"
)

type RouteGroup struct {
	router      *Router
	prefix      string
	middlewares []routing.Handler
}

func (r *Router) Get(uri string, handler routing.Handler) {
	r.New(uri, handler, MethodGet, nil)
}

func (r *Router) Post(uri string, handler routing.Handler) {
	r.New(uri, handler, MethodPost, nil)
}

func (r *Router) Put(uri string, handler routing.Handler) {
	r.New(uri, handler, MethodPut, nil)
}

func (r *Router) Delete(uri string, handler routing.Handler) {
	r.New(uri, handler, MethodDelete, nil)
}

func (r *Router) Handle(uri string, methods []string, handler routing.Handler, middlewares ...routing.Handler) {
	r.New(uri, handler, methods, middlewares)
}

func (r *Router) Group(prefix string, middlewares ...routing.Handler) *RouteGroup {
	return &RouteGroup{
		router:      r,
		prefix:      prefix,
		middlewares: append([]routing.Handler{}, middlewares...),
	}
}

func (r *Router) Version(prefix string, middlewares ...routing.Handler) *RouteGroup {
	return r.Group(prefix, middlewares...)
}

func (g *RouteGroup) Handle(uri string, methods []string, handler routing.Handler, middlewares ...routing.Handler) {
	all := append([]routing.Handler{}, g.middlewares...)
	all = append(all, middlewares...)
	g.router.New(joinRoute(g.prefix, uri), handler, methods, all)
}

func (g *RouteGroup) Get(uri string, handler routing.Handler, middlewares ...routing.Handler) {
	g.Handle(uri, MethodGet, handler, middlewares...)
}

func (g *RouteGroup) Post(uri string, handler routing.Handler, middlewares ...routing.Handler) {
	g.Handle(uri, MethodPost, handler, middlewares...)
}

func (g *RouteGroup) Put(uri string, handler routing.Handler, middlewares ...routing.Handler) {
	g.Handle(uri, MethodPut, handler, middlewares...)
}

func (g *RouteGroup) Delete(uri string, handler routing.Handler, middlewares ...routing.Handler) {
	g.Handle(uri, MethodDelete, handler, middlewares...)
}

func (r *Router) New(uri string, handler routing.Handler, methods []string, middlewares []routing.Handler) {
	for _, method := range methods {
		if r.app != nil {
			r.app.logInfo(fmt.Sprintf("register %s %s", method, uri))
		}

		r.Routes = append(r.Routes, route{
			Method:      method,
			URI:         uri,
			Handler:     handler,
			Middlewares: middlewares,
		})
	}
}

func (r *Router) load() {
	for _, route := range r.Routes {
		handlers := attachMiddleware(r.app, route.Middlewares, route.Handler)

		switch route.Method {
		case "GET":
			r.Router.Get(route.URI, handlers...)
		case "PUT":
			r.Router.Put(route.URI, handlers...)
		case "POST":
			r.Router.Post(route.URI, handlers...)
		case "PATCH":
			r.Router.Patch(route.URI, handlers...)
		case "DELETE":
			r.Router.Delete(route.URI, handlers...)
		case "OPTIONS":
			r.Router.Options(route.URI, handlers...)
		}
	}
}

func attachMiddleware(cfg *Config, routeMiddlewares []routing.Handler, routeHandler func(ctx *Context) error) []routing.Handler {
	var middleware []routing.Handler

	if cfg != nil {
		middleware = append(middleware, injectConfig(cfg))

		if cfg.Middleware.Limiter.Enabled {
			middleware = append(middleware, Limiter)
		}
		if cfg.Middleware.LogRequest.Enabled {
			middleware = append(middleware, LogRequest)
		}
		if cfg.Middleware.BasicHeaders.Enabled {
			middleware = append(middleware, BasicHeaders)
		}
		if cfg.Middleware.AdminOnly.Enabled {
			middleware = append(middleware, AdminOnly)
		}
		if cfg.Middleware.BasicAuth.Enabled {
			middleware = append(middleware, BasicAuth)
		}
	}

	// Route-specific middleware should run after framework-global middleware
	// and before the final route handler.
	if len(routeMiddlewares) > 0 {
		middleware = append(middleware, routeMiddlewares...)
	}

	// our route handler should be after all middleware
	middleware = append(middleware, routeHandler)
	return middleware
}

// wireDefaultHandlers attempts to explicitly configure router-level default handlers.
// It is intentionally reflection-based to avoid compile regressions across fasthttp-routing variants.
func wireDefaultHandlers(r *routing.Router) {
	if r == nil {
		return
	}

	set := func(name string, h routing.Handler) {
		// Prefer method: r.NotFound(handler) / r.MethodNotAllowed(handler) if present
		mv := reflect.ValueOf(r).MethodByName(name)
		if mv.IsValid() && mv.Type().NumIn() == 1 {
			arg := reflect.ValueOf(h)
			in0 := mv.Type().In(0)
			if arg.Type().AssignableTo(in0) {
				mv.Call([]reflect.Value{arg})
				return
			}
			if arg.Type().ConvertibleTo(in0) {
				mv.Call([]reflect.Value{arg.Convert(in0)})
				return
			}
		}

		// Fallback: settable field NotFound / MethodNotAllowed if present
		ev := reflect.ValueOf(r)
		if ev.Kind() == reflect.Ptr {
			ev = ev.Elem()
		}
		if !ev.IsValid() {
			return
		}

		fv := ev.FieldByName(name)
		if fv.IsValid() && fv.CanSet() {
			hv := reflect.ValueOf(h)
			if hv.Type().AssignableTo(fv.Type()) {
				fv.Set(hv)
				return
			}
			if hv.Type().ConvertibleTo(fv.Type()) {
				fv.Set(hv.Convert(fv.Type()))
				return
			}
		}
	}

	set("NotFound", NotFoundServer)
	set("MethodNotAllowed", MethodNotAllowed)
	set("MethodNotAllowedHandler", MethodNotAllowed)
}
