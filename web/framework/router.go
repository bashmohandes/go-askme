package framework

import (
	"log"
	"net/http"
)

// Router interface
type Router interface {
	Routes() []*Route
	Get(string, RouteHandler) *Route
	Post(string, RouteHandler) *Route
	Delete(string, RouteHandler) *Route
	Put(string, RouteHandler) *Route
}

// RouteHandler defines handler
type RouteHandler func(context Context)

// FileProvider fetches files by name
type FileProvider interface {
	List() []string
	String(name string) string
	Open(name string) (http.File, error)
}

// Route captures the http actions of the controller
type Route struct {
	Method  string
	Path    string
	Options *RouteOptions
	Func    RouteHandler
}

// RouteOptions settings
type RouteOptions struct {
	AuthRequired bool
}

// router represents router in MVC  model
type router struct {
	actions []*Route
}

// NewRouter initializes the controller
func NewRouter() Router {
	return &router{
		actions: make([]*Route, 0),
	}
}

// Routes return the list of configured actions
func (r *router) Routes() []*Route {
	return r.actions
}

func (r *router) Get(path string, f RouteHandler) *Route {
	return r.route("GET", path, f)
}

func (r *router) Post(path string, f RouteHandler) *Route {
	return r.route("POST", path, f)
}

func (r *router) Delete(path string, f RouteHandler) *Route {
	return r.route("DELETE", path, f)
}

func (r *router) Put(path string, f RouteHandler) *Route {
	return r.route("PUT", path, f)
}

// AddAction adds route to controller
func (r *router) route(method string, path string, f RouteHandler) *Route {
	log.Printf("caller %T, method %s, path %s\n", r, method, path)
	route := &Route{Method: method, Path: path, Func: f, Options: &RouteOptions{}}
	r.actions = append(r.actions, route)
	return route
}

// Authenticated Adds authentication options
func (rt *Route) Authenticated() *Route {
	rt.Options.AuthRequired = true
	return rt
}
