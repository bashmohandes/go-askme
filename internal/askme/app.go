package askme

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bashmohandes/go-askme/internal/askme/controllers"
	"github.com/bashmohandes/go-askme/internal/service"
	"github.com/bashmohandes/go-askme/internal/shared"
)

// Server represents the AskMe application server
type Server struct {
	config          *Config
	fileProvider    common.FileProvider
	questionService service.QuestionService
}

// Config configuration
type Config struct {
	// Assets relative path to askme package
	Assets string
	Port   int
}

//Start method starts the AskMe App
func (server *Server) Start() {
	b := controllers.Blog(server.questionService, server.fileProvider)
	fmt.Println("Hello!")
	fmt.Printf("Listening on port %d\n", server.config.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", server.config.Port), b))
}

// NewServer Creates a new AskMe app server
func NewServer(
	config *Config,
	fileProvider common.FileProvider,
	questionService service.QuestionService) *Server {
	return &Server{
		config:          config,
		fileProvider:    fileProvider,
		questionService: questionService,
	}
}
