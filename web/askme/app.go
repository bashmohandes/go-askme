package askme

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bashmohandes/go-askme/shared"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/askme/controllers"
)

// Server represents the AskMe application server
type Server struct {
	config       *Config
	fileProvider shared.FileProvider
	asksScenario user.AsksUsecase
}

// Config configuration
type Config struct {
	// Assets relative path to askme package
	Assets string
	Port   int
}

//Start method starts the AskMe App
func (server *Server) Start() {
	b := controllers.Blog(server.asksScenario, server.fileProvider)
	fmt.Println("Hello!")
	fmt.Printf("Listening on port %d\n", server.config.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", server.config.Port), b))
}

// NewServer Creates a new AskMe app server
func NewServer(
	config *Config,
	fileProvider shared.FileProvider,
	scenario user.AsksUsecase) *Server {
	return &Server{
		config:       config,
		fileProvider: fileProvider,
		asksScenario: scenario,
	}
}
