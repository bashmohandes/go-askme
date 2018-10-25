package framework

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// App represents the AskMe application server
type App interface {
	Start() error
	Use(m Middleware)
	UseFunc(f MiddlewareFunc)
}

type app struct {
	config       *Config
	fileProvider FileProvider
	rtr          Router
	middlewares  []Middleware
	sessionStore SessionManager
}

// Config configuration
type Config struct {
	Debug        bool
	Port         int
	PublicFolder string
}

//Start method starts the AskMe App
func (app *app) Start() error {
	mux := httprouter.New()

	for _, r := range app.rtr.Routes() {
		mux.Handle(r.Method, r.Path, app.handle(r.Func, r))
	}

	mux.ServeFiles("/public/*filepath", app.fileProvider)
	fmt.Printf("Listening on port %d\n", app.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.Port), mux)
}

func (app *app) Use(m Middleware) {
	app.middlewares = append(app.middlewares, m)
}

func (app *app) UseFunc(f MiddlewareFunc) {
	app.middlewares = append(app.middlewares, f)
}

func (app *app) handle(handler RouteHandler, route *Route) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		c := app.initContext(w, r, params)
		for _, m := range app.middlewares {
			success := m.Run(c)
			if !success {
				return
			}
		}
		if route.Options.AuthRequired && c.User() == nil {
			c.Redirect(fmt.Sprintf("/login?redir=%s", r.RequestURI), http.StatusTemporaryRedirect)
		}
		handler(c)
	}
}

func (app *app) initContext(w http.ResponseWriter, r *http.Request, p httprouter.Params) Context {
	c := &cxt{w: w, r: r, p: p}
	session := app.sessionStore.FetchOrCreate(c)
	c.s = session
	return c
}

// NewApp Creates a new app server
func NewApp(
	config *Config,
	ctrl Router,
	fileProvider FileProvider,
	sessionStore SessionManager) App {
	return &app{
		config:       config,
		fileProvider: fileProvider,
		rtr:          ctrl,
		middlewares:  make([]Middleware, 0),
		sessionStore: sessionStore,
	}
}
