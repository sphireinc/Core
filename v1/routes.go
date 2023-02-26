package core

import (
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
)

// Get is a shorthand for New() specifically for the GET method
func (r *Router) Get(uri string, handler routing.Handler) {
	r.New(uri, handler, MethodGet, nil)
}

// Post is a shorthand for New() specifically for the POST method
func (r *Router) Post(uri string, handler routing.Handler) {
	r.New(uri, handler, MethodPost, nil)
}

// Put is a shorthand for New() specifically for the PUT method
func (r *Router) Put(uri string, handler routing.Handler) {
	r.New(uri, handler, MethodPut, nil)
}

// Delete is a shorthand for New() specifically for the DELETE method
func (r *Router) Delete(uri string, handler routing.Handler) {
	r.New(uri, handler, MethodDelete, nil)
}

// New appends a new route to our route struct.
func (r *Router) New(uri string, handler routing.Handler, methods []string, middlewares []routing.Handler) {
	for _, method := range methods {
		App.Log.Writer.Info(fmt.Sprintf("register %s %s", method, uri))
		r.Routes = append(r.Routes, route{
			Method:      method,
			URI:         uri,
			Handler:     handler,
			Middlewares: middlewares,
		})
	}
}

// load all of our routes
func (r *Router) load() {
	for _, route := range r.Routes {
		switch route.Method {
		case "GET":
			r.Router.Get(route.URI, attachMiddleware(route.Handler)...)
		case "PUT":
			r.Router.Put(route.URI, attachMiddleware(route.Handler)...)
		case "POST":
			r.Router.Post(route.URI, attachMiddleware(route.Handler)...)
		case "PATCH":
			r.Router.Patch(route.URI, attachMiddleware(route.Handler)...)
		case "DELETE":
			r.Router.Delete(route.URI, attachMiddleware(route.Handler)...)
		case "OPTIONS":
			r.Router.Options(route.URI, attachMiddleware(route.Handler)...)
		}
	}
}

func attachMiddleware(routeHandler func(ctx *Context) error) []routing.Handler {
	var middleware []routing.Handler

	if App.Middleware.Limiter.Enabled {
		middleware = append(middleware, Limiter)
	}
	if App.Middleware.LogRequest.Enabled {
		middleware = append(middleware, LogRequest)
	}
	if App.Middleware.BasicHeaders.Enabled {
		middleware = append(middleware, BasicHeaders)
	}
	if App.Middleware.AdminOnly.Enabled {
		middleware = append(middleware, AdminOnly)
	}
	if App.Middleware.BasicAuth.Enabled {
		middleware = append(middleware, BasicAuth)
	}

	// our route handler should be after all internal middleware
	middleware = append(middleware, routeHandler)
	return middleware
}
