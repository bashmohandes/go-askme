package question

import "github.com/bashmohandes/go-askme/model"

// Repository interface
type Repository interface {
	LoadUnansweredQuestions(user uint) ([]*models.Question, error)
	Add(question *models.Question) (*models.Question, error)
	GetByID(id uint) (*models.Question, error)
}
