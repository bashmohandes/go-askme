package main

import (
	"github.com/bashmohandes/go-askme/internal/askme"
	"github.com/bashmohandes/go-askme/internal/data"
	"github.com/bashmohandes/go-askme/internal/domain"
	"github.com/gobuffalo/packr"
)

func main() {
	config := &askme.Config{
		Assets: "../../internal/askme/public",
		Port:   8080,
	}
	box := packr.NewBox(config.Assets)
	models.QuestionsRepo = repository.NewQuestionRepo()
	models.AnswersRepo = repository.NewAnswerRepository()
	appServer := askme.NewServer(config, box)

	appServer.Start()
}
