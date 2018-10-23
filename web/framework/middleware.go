package framework

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Middleware interface
type Middleware interface {
	Run(http.ResponseWriter, *http.Request, httprouter.Params) bool
}

// MiddlewareFunc adapter for funcs to middleware
type MiddlewareFunc func(http.ResponseWriter, *http.Request, httprouter.Params) bool

// Run middleware func
func (f MiddlewareFunc) Run(w http.ResponseWriter, r *http.Request, p httprouter.Params) bool {
	return f(w, r, p)
}
