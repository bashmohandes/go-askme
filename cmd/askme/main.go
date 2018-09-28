package main

import (
	"github.com/bashmohandes/go-askme/internal/askme"
	"github.com/bashmohandes/go-askme/internal/data"
	"github.com/bashmohandes/go-askme/internal/service"
	"github.com/bashmohandes/go-askme/internal/shared"
	"github.com/gobuffalo/packr"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	container.Provide(func() *askme.Config {
		return &askme.Config{
			Assets: "../../internal/askme/public",
			Port:   8080,
		}
	})
	container.Provide(func(config *askme.Config) common.FileProvider {
		return packr.NewBox(config.Assets)
	})
	container.Provide(askme.NewServer)
	container.Provide(repository.NewQuestionRepository)
	container.Provide(repository.NewAnswerRepository)
	container.Provide(service.NewQuestionService)
	err := container.Invoke(func(server *askme.Server) {
		server.Start()
	})

	if err != nil {
		panic(err)
	}
}
