package core

import (
	"fmt"
	routing "github.com/qiangxue/fasthttp-routing"
)

// load all of our routes
func (r *Router) load() {
	for _, route := range r.Routes {
		switch route.Method {
		case "GET":
			r.Router.Get(route.URI, LogRequest, BasicHeaders, AdminOnly, BasicAuth, Limiter, route.Handler)
		case "PUT":
			r.Router.Put(route.URI, LogRequest, BasicHeaders, AdminOnly, BasicAuth, Limiter, route.Handler)
		case "POST":
			r.Router.Post(route.URI, LogRequest, BasicHeaders, AdminOnly, BasicAuth, Limiter, route.Handler)
		case "PATCH":
			r.Router.Patch(route.URI, LogRequest, BasicHeaders, AdminOnly, BasicAuth, Limiter, route.Handler)
		case "DELETE":
			r.Router.Delete(route.URI, LogRequest, BasicHeaders, AdminOnly, BasicAuth, Limiter, route.Handler)
		case "OPTIONS":
			r.Router.Options(route.URI, LogRequest, BasicHeaders, AdminOnly, BasicAuth, Limiter, route.Handler)
		}
	}
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
