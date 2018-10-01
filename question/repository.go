package question

import "github.com/bashmohandes/go-askme/model"

// Repository interface
type Repository interface {
	LoadQuestions(user models.UniqueID) []*models.Question
	Save(question *models.Question) *models.Question
}
