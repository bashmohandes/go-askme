package askme

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bashmohandes/go-askme/shared"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/askme/controllers"
	"github.com/julienschmidt/httprouter"
)

// Server represents the AskMe application server
type Server struct {
	config        *Config
	fileProvider  shared.FileProvider
	asksScenario  user.AsksUsecase
	answrScenario user.AnswersUsecase
}

// Config configuration
type Config struct {
	// Assets relative path to askme package
	Assets string
	Port   int
}

//Start method starts the AskMe App
func (server *Server) Start() {
	c := controllers.New(server.asksScenario, server.answrScenario, server.fileProvider)
	mux := httprouter.New()

	for _, a := range c.Actions() {
		mux.Handle(a.Method, a.Path, a.Func)
	}

	mux.ServeFiles("/public/*filepath", server.fileProvider)

	fmt.Println("Hello!")
	fmt.Printf("Listening on port %d\n", server.config.Port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", server.config.Port), mux))
}

// NewServer Creates a new AskMe app server
func NewServer(
	config *Config,
	fileProvider shared.FileProvider,
	asksUC user.AsksUsecase,
	answrUC user.AnswersUsecase) *Server {
	return &Server{
		config:        config,
		fileProvider:  fileProvider,
		asksScenario:  asksUC,
		answrScenario: answrUC,
	}
}
