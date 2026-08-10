package main

import (
	"fmt"
	"log"

	"github.com/bashmohandes/go-askme/answer/db"
	"github.com/bashmohandes/go-askme/models"
	"github.com/bashmohandes/go-askme/question/db"
	userRepo "github.com/bashmohandes/go-askme/user/db"
	"github.com/bashmohandes/go-askme/user/usecase"
	"github.com/bashmohandes/go-askme/web/askme"
	"github.com/bashmohandes/go-askme/web/askme/controllers"
	"github.com/bashmohandes/go-askme/web/framework"
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	container.Provide(framework.NewConnection)
	container.Provide(question.NewRepository)
	container.Provide(answer.NewRepository)
	container.Provide(userRepo.NewRepository)
	container.Provide(user.NewAsksUsecase)
	container.Provide(user.NewAnswersUsecase)
	container.Provide(user.NewAuthUsecase)
	container.Provide(controllers.NewHomeController)
	container.Provide(controllers.NewProfileController)
	container.Provide(controllers.NewOktaController)
	container.Provide(askme.NewApp)
	err := container.Invoke(func(app *askme.App, config *framework.Config) {
		err := migrateDB(config)
		if err != nil {
			log.Fatal(err)
		}
		if e := app.Start(); e != nil {
			log.Fatalln(e)
		}
	})

	if err != nil {
		log.Fatal(err)
	}
}

func newFileProvider(config *framework.Config) framework.FileProvider {
	return packr.NewBox(config.PublicFolder)
}

func newConfig() (*framework.Config, error) {
	return framework.LoadConfig()
}

func migrateDB(config *framework.Config) error {
	log.Print("Auto Migration Starting")
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresHost, config.PostgresPort, config.PostgresDB)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("open database for migration: %w", err)
	}
	defer db.Close()
	db.LogMode(config.Debug)
	if err := db.AutoMigrate(&models.User{}, &models.Question{}, &models.Answer{}).Error; err != nil {
		return fmt.Errorf("auto migrate database: %w", err)
	}
	log.Print("Auto Migration Ended")
	return nil
}
