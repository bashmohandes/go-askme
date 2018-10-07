package main

import (
	"github.com/bashmohandes/go-askme/answer/inmemory"
	"github.com/bashmohandes/go-askme/question/inmemory"
	"github.com/bashmohandes/go-askme/shared"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/askme"
	"github.com/gobuffalo/packr"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()
	container.Provide(newConfig)
	container.Provide(newFileProvider)
	container.Provide(askme.NewServer)
	container.Provide(question.NewRepository)
	container.Provide(answer.NewRepository)
	container.Provide(user.NewAsksUsecase)
	container.Provide(user.NewAnswersUsecase)
	err := container.Invoke(func(server *askme.Server) {
		server.Start()
	})

	if err != nil {
		panic(err)
	}
}

func newFileProvider(config *askme.Config) shared.FileProvider {
	return packr.NewBox(config.Assets)
}

func newConfig() *askme.Config {
	return &askme.Config{
		Assets: "../../web/askme/public",
		Port:   8080,
	}
}
