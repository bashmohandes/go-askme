package service

import "github.com/bashmohandes/go-askme/internal/domain"

// QuestionService type
type QuestionService interface {
	LoadQuestions(userID models.UniqueID) []*models.Question
	Save(question *models.Question) *models.Question
}
