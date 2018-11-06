package question

import "github.com/bashmohandes/go-askme/model"

// Repository interface
type Repository interface {
	LoadUnansweredQuestions(user uint) []*models.Question
	Add(question *models.Question)
}
