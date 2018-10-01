package service

import (
	"github.com/bashmohandes/go-askme/internal/domain"
	"github.com/bashmohandes/go-askme/internal/repository"
	"github.com/bashmohandes/go-askme/internal/service"
)

type questionService struct {
	questionRepo repository.QuestionRepository
	answerRepo   repository.AnswerRepository
}

// NewQuestionService creates a new service
func NewQuestionService(qRepo repository.QuestionRepository, aRepo repository.AnswerRepository) service.QuestionService {
	return &questionService{
		questionRepo: qRepo,
		answerRepo:   aRepo,
	}
}

// LoadQuestions load questions model
func (svc *questionService) LoadQuestions(userID models.UniqueID) []*models.Question {
	return svc.questionRepo.LoadQuestions(userID)
}

// Saves question
func (svc *questionService) Save(question *models.Question) *models.Question {
	return svc.questionRepo.Save(question)
}
