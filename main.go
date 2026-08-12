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
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config, err := framework.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	app, err := buildApplication(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := migrateDB(config); err != nil {
		log.Fatal(err)
	}
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func buildApplication(config *framework.Config) (*askme.App, error) {
	publicFiles, err := askme.NewPublicFiles()
	if err != nil {
		return nil, fmt.Errorf("load embedded public files: %w", err)
	}

	router := framework.NewRouter()
	renderer := framework.NewRenderer(publicFiles, config)
	sessions := framework.NewInMemorySessionStore(config)
	connection := framework.NewConnection(config)
	questionRepository := question.NewRepository(connection)
	answerRepository := answer.NewRepository(connection)
	userRepository := userRepo.NewRepository(connection)
	asks := user.NewAsksUsecase(questionRepository, answerRepository, userRepository)
	answers := user.NewAnswersUsecase(questionRepository, answerRepository, userRepository)
	auth := user.NewAuthUsecase(userRepository)
	homeController := controllers.NewHomeController(router)
	profileController := controllers.NewProfileController(router, renderer, asks, answers)
	oktaController := controllers.NewOktaController(router, renderer, config, sessions, auth)
	base := framework.NewApp(config, router, publicFiles, sessions)

	return askme.NewApp(base, homeController, profileController, oktaController), nil
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
