package framework

// Middleware interface
type Middleware interface {
	Run(context Context) bool
}

// MiddlewareFunc adapter for funcs to middleware
type MiddlewareFunc func(context Context) bool

// Run middleware func
func (f MiddlewareFunc) Run(context Context) bool {
	return f(context)
}
