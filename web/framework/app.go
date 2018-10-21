package framework

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// App represents the AskMe application server
type App struct {
	config       *Config
	fileProvider FileProvider
	rtr          Router
}

// Config configuration
type Config struct {
	Debug        bool
	Port         int
	PublicFolder string
}

//Start method starts the AskMe App
func (app *App) Start() error {
	mux := httprouter.New()

	for _, a := range app.rtr.Actions() {
		log.Printf("Method %s, Path %s\n", a.Method, a.Path)
		mux.Handle(a.Method, a.Path, a.Func)
	}

	mux.ServeFiles("/public/*filepath", app.fileProvider)

	fmt.Println("Hello!")
	fmt.Printf("Listening on port %d\n", app.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.Port), mux)
}

// NewApp Creates a new app server
func NewApp(
	config *Config,
	ctrl Router, fileProvider FileProvider) *App {
	return &App{
		config:       config,
		fileProvider: fileProvider,
		rtr:          ctrl,
	}
}
