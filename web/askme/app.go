package askme

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bashmohandes/go-askme/web"

	"github.com/bashmohandes/go-askme/shared"
	"github.com/bashmohandes/go-askme/web/askme/controllers"
	"github.com/julienschmidt/httprouter"
)

// Server represents the AskMe application server
type Server struct {
	config            *framework.Config
	fileProvider      shared.FileProvider
	homeController    *controllers.HomeController
	profileController *controllers.ProfileController
}

//Start method starts the AskMe App
func (server *Server) Start() {
	mux := httprouter.New()

	actions := make([]*framework.Action, 0)

	actions = append(actions, server.homeController.Actions()...)
	actions = append(actions, server.profileController.Actions()...)

	for _, a := range actions {
		mux.Handle(a.Method, a.Path, a.Func)
	}

	mux.ServeFiles("/public/*filepath", server.fileProvider)

	fmt.Println("Hello!")
	fmt.Printf("Listening on port %d\n", server.config.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", server.config.Port), mux))
}

// NewServer Creates a new AskMe app server
func NewServer(
	config *framework.Config,
	fileProvider shared.FileProvider,
	hc *controllers.HomeController,
	pc *controllers.ProfileController) *Server {
	return &Server{
		config:            config,
		fileProvider:      fileProvider,
		homeController:    hc,
		profileController: pc,
	}
}
