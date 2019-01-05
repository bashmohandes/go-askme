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

type RouteAdapter func(f RouteHandler) RouteHandler

func AdaptRoute(h RouteHandler, adapters ...RouteAdapter) RouteHandler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}