package service

import (
	"github.com/bashmohandes/go-askme/internal/domain"
	"github.com/bashmohandes/go-askme/internal/repository"
)

// QuestionService type
type QuestionService struct {
	questionRepo repository.QuestionRepository
}

// NewQuestionService creates a new service
func NewQuestionService(repo repository.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepo: repo,
	}
}

// LoadQuestions load questions model
func (svc *QuestionService) LoadQuestions(userID models.UniqueID) []*models.Question {
	return svc.questionRepo.LoadQuestions(userID)
}
