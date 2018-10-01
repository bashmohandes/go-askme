package usecase

import (
	"github.com/bashmohandes/go-askme/answer"
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/question"
	"github.com/bashmohandes/go-askme/user"
)

type userUsecase struct {
	questionRepo question.Repository
	answerRepo   answer.Repository
}

// NewUsecase creates a new service
func NewUsecase(qRepo question.Repository, aRepo answer.Repository) user.Usecase {
	return &userUsecase{
		questionRepo: qRepo,
		answerRepo:   aRepo,
	}
}

// LoadQuestions load questions model
func (svc *userUsecase) FetchUnansweredQuestions(userID models.UniqueID) []*models.Question {
	return svc.questionRepo.LoadQuestions(userID)
}

// Saves question
func (svc *userUsecase) Ask(question *models.Question) *models.Question {
	return svc.questionRepo.Save(question)
}
