package user

import "github.com/bashmohandes/go-askme/model"

// Usecase type
type Usecase interface {
	FetchUnansweredQuestions(userID models.UniqueID) []*models.Question
	Ask(from *models.User, to *models.User, question string) *models.Question
	Like(user *models.User, answer *models.Answer) uint
}
