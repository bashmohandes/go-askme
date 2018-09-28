package main

import (
	"fmt"

	"github.com/bashmohandes/go-askme/internal/askme"
	"github.com/bashmohandes/go-askme/internal/data"
	"github.com/bashmohandes/go-askme/internal/domain"
)

func main() {
	fmt.Println("Hello!")
	models.QuestionsRepo = repository.NewQuestionRepo()
	models.AnswersRepo = repository.NewAnswerRepository()
	askme.Start()
}
