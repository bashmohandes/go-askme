package main

import (
	"github.com/bashmohandes/go-askme/internal/askme"
	"github.com/bashmohandes/go-askme/internal/repository/inmemory"
	"github.com/bashmohandes/go-askme/internal/service/default"
	"github.com/bashmohandes/go-askme/internal/shared"
	"github.com/gobuffalo/packr"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()
	container.Provide(newConfig)
	container.Provide(newFileProvider)
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

func newFileProvider(config *askme.Config) common.FileProvider {
	return packr.NewBox(config.Assets)
}

func newConfig() *askme.Config {
	return &askme.Config{
		Assets: "../../internal/askme/public",
		Port:   8080,
	}
}
