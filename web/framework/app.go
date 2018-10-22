package framework

import (
	"fmt"
	"log"
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

	for _, a := range app.rtr.Actions() {
		log.Printf("Method %s, Path %s\n", a.Method, a.Path)
		mux.Handle(a.Method, a.Path, app.handle(a.Func))
	}

	mux.ServeFiles("/public/*filepath", app.fileProvider)

	fmt.Println("Hello!")
	fmt.Printf("Listening on port %d\n", app.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.Port), mux)
}

func (app *app) Use(m Middleware) {
	app.middlewares = append(app.middlewares, m)
}

func (app *app) UseFunc(f MiddlewareFunc) {
	app.middlewares = append(app.middlewares, f)
}

func (app *app) handle(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		for _, m := range app.middlewares {
			success := m.Run(w, r, params)
			if !success {
				return
			}
		}

		handler(w, r, params)
	}
}

// NewApp Creates a new app server
func NewApp(
	config *Config,
	ctrl Router,
	fileProvider FileProvider) App {
	return &app{
		config:       config,
		fileProvider: fileProvider,
		rtr:          ctrl,
		middlewares:  make([]Middleware, 0),
	}
}
