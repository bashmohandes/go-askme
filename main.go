package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	answer "github.com/bashmohandes/go-askme/answer/db"
	"github.com/bashmohandes/go-askme/models"
	question "github.com/bashmohandes/go-askme/question/db"
	userRepo "github.com/bashmohandes/go-askme/user/db"
	user "github.com/bashmohandes/go-askme/user/usecase"
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
	container.Provide(controllers.NewAuthController)
	container.Provide(askme.NewApp)
	err := container.Invoke(func(app *askme.App) {
		err := migrateDB()
		if err != nil {
			log.Fatal(err)
		}
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

	config := &framework.Config{
		Debug:              debug,
		PublicFolder:       os.Getenv("PUBLIC_FOLDER"),
		Port:               port,
		SessionMaxLifeTime: sessionMaxLife,
		SessionCookie:      os.Getenv("SESSION_COOKIE"),
		PostgresUser:       os.Getenv("POSTGRES_USER"),
		PostgresPassword:   os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:         os.Getenv("POSTGRES_DB"),
		PostgresHost:       os.Getenv("POSTGRES_HOST"),
		OktaClient:         os.Getenv("OKTA_CLIENT_ID"),
		OktaSecret:         os.Getenv("OKTA_CLIENT_SECRET"),
		OktaIssuer:         os.Getenv("OKTA_ISSUER"),
	}

	linkedInIdp := os.Getenv("OKTA_SOCIAL_LINKEDIN_IDP")
	facebookIdp := os.Getenv("OKTA_SOCIAL_FACEBOOK_IDP")

	if linkedInIdp != "" {
		config.OktaSocialIdps = append(config.OktaSocialIdps, framework.OktaSocialIdp{
			ID:   linkedInIdp,
			Name: "LINKEDIN",
		})
	}

	if facebookIdp != "" {
		config.OktaSocialIdps = append(config.OktaSocialIdps, framework.OktaSocialIdp{
			ID:   facebookIdp,
			Name: "FACEBOOK",
		})
	}

	return config
}

func migrateDB() error {
	log.Print("Auto Migration Starting")
	config := newConfig()
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresHost, 5432, config.PostgresDB)
	db, err := gorm.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatalf("%v", err)
	}
	db.LogMode(config.Debug)
	db.AutoMigrate(&models.User{}, &models.Question{}, &models.Answer{})
	log.Print("Auto Migration Ended")
	return nil
}
