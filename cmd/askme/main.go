package main

import (
	"log"
	"os"
	"strconv"

	"github.com/bashmohandes/go-askme/answer/inmemory"
	"github.com/bashmohandes/go-askme/cmd/askme/providers"
	"github.com/bashmohandes/go-askme/question/inmemory"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/askme"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()
	container.Provide(newConfig)
	container.Provide(providers.NewFileProvider)
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

func newConfig() *askme.Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Incorrect format: %v\n", err)
	}
	reload, err := strconv.ParseBool(os.Getenv("RELOAD_ASSETS_ON_CHANGE"))
	if err != nil {
		log.Fatalf("Incorrect format: %v\n", err)
	}

	return &askme.Config{
		Assets:               "../../../web/askme/public",
		Port:                 port,
		ReloadAssetsOnChange: reload,
	}
}
