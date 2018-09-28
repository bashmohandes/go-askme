package service

import "github.com/bashmohandes/go-askme/internal/domain"

// QuestionService type
type QuestionService struct {
	questionRepo models.QuestionRepository
}

// NewQuestionService creates a new service
func NewQuestionService(repo models.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepo: repo,
	}
}

// LoadQuestions load questions model
func (svc *QuestionService) LoadQuestions(userID models.UniqueID) []models.Question {
	return svc.questionRepo.LoadQuestions(userID)
}
