package main

import (
	answerRepo "github.com/bashmohandes/go-askme/answer/repository"
	questionRepo "github.com/bashmohandes/go-askme/question/repository"
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
	container.Provide(questionRepo.NewRepository)
	container.Provide(answerRepo.NewRepository)
	container.Provide(usecase.NewUsecase)
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
		Assets: "../../web/askme/public",
		Port:   8080,
	}
}
