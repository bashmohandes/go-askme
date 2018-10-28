package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bashmohandes/go-askme/answer/inmemory"
	"github.com/bashmohandes/go-askme/question/inmemory"
	userRepo "github.com/bashmohandes/go-askme/user/inmemory"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/askme"
	"github.com/bashmohandes/go-askme/web/askme/controllers"
	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/gobuffalo/packr"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()
	container.Provide(newConfig)
	container.Provide(newFileProvider)
	container.Provide(framework.NewApp)
	container.Provide(framework.NewRouter)
	container.Provide(framework.NewRenderer)
	container.Provide(framework.NewInMemorySessionStore)
	container.Provide(question.NewRepository)
	container.Provide(answer.NewRepository)
	container.Provide(userRepo.NewRepository)
	container.Provide(user.NewAsksUsecase)
	container.Provide(user.NewAnswersUsecase)
	container.Provide(user.NewAuthUsecase)
	container.Provide(controllers.NewHomeController)
	container.Provide(controllers.NewProfileController)
	container.Provide(controllers.NewAuthController)
	container.Provide(askme.NewApp)
	err := container.Invoke(func(app *askme.App) {
		if e := app.Start(); e != nil {
			log.Fatalln(e)
		}
	})

	if err != nil {
		panic(err)
	}
}

func newFileProvider(config *framework.Config) framework.FileProvider {
	return packr.NewBox(config.PublicFolder)
}

func newConfig() *framework.Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Incorrect 'PORT' format: %v\n", err)
	}
	debug, err := strconv.ParseBool(os.Getenv("DEBUG_MODE"))
	if err != nil {
		log.Fatalf("Incorrect 'DEBUG_MODE' format: %v\n", err)
	}
	sessionMaxLife, err := time.ParseDuration(os.Getenv("SESSION_MAX_LIFE_TIME"))
	if err != nil {
		log.Fatalf("Incorrect 'SESSION_MAX_LIFE_TIME' format: %v\n", err)
	}
	sessionCookie := os.Getenv("SESSION_COOKIE")
	public := os.Getenv("PUBLIC_FOLDER")
	return &framework.Config{
		Debug:              debug,
		PublicFolder:       public,
		Port:               port,
		SessionMaxLifeTime: sessionMaxLife,
		SessionCookie:      sessionCookie,
	}
}
