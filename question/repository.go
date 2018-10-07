package question

import "github.com/bashmohandes/go-askme/model"

// Repository interface
type Repository interface {
	LoadUnansweredQuestions(user models.UniqueID) []*models.Question
	Add(question *models.Question)
}
