package repository

import "github.com/bashmohandes/go-askme/internal/domain"

// QuestionRepository interface
type QuestionRepository interface {
	LoadQuestions(user models.UniqueID) []*models.Question
	Save(question *models.Question) *models.Question
}
