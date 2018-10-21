package framework

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Router interface
type Router interface {
	Actions() []*Action
	Get(path string, f httprouter.Handle)
	Post(path string, f httprouter.Handle)
	Delete(path string, f httprouter.Handle)
	Put(path string, f httprouter.Handle)
}

// FileProvider fetches files by name
type FileProvider interface {
	List() []string
	String(name string) string
	Open(name string) (http.File, error)
}

// Action captures the http actions of the controller
type Action struct {
	Method string
	Path   string
	Func   httprouter.Handle
}

// router represents router in MVC  model
type router struct {
	actions []*Action
}

// NewRouter initializes the controller
func NewRouter() Router {
	return &router{
		actions: make([]*Action, 0),
	}
}

// Actions return the list of configured actions
func (r *router) Actions() []*Action {
	return r.actions
}

func (r *router) Get(path string, f httprouter.Handle) {
	r.action("GET", path, f)
}

func (r *router) Post(path string, f httprouter.Handle) {
	r.action("POST", path, f)
}

func (r *router) Delete(path string, f httprouter.Handle) {
	r.action("DELETE", path, f)
}

func (r *router) Put(path string, f httprouter.Handle) {
	r.action("PUT", path, f)
}

// AddAction adds action to controller
func (r *router) action(method string, path string, f httprouter.Handle) {
	log.Printf("caller %T, method %s, path %s\n", r, method, path)
	r.actions = append(r.actions, &Action{method, path, f})
}
