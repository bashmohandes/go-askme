package user

import "github.com/bashmohandes/go-askme/model"

// Usecase type
type Usecase interface {
	FetchUnansweredQuestions(userID models.UniqueID) []*models.Question
	Ask(question *models.Question) *models.Question
}
