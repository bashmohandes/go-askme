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
	return svc.questionRepo.LoadUnansweredQuestions(userID)
}

// Saves question
func (svc *userUsecase) Ask(from *models.User, to *models.User, question string) *models.Question {
	q := from.Ask(to, question)
	return svc.questionRepo.Save(q)
}

// Likes a question
func (svc *userUsecase) Like(user *models.User, answer *models.Answer) uint {
	svc.answerRepo.AddLike(answer, user)
	return svc.answerRepo.GetLikesCount(answer)
}

// Unlikes a question
func (svc *userUsecase) Unlike(user *models.User, answer *models.Answer) uint {
	svc.answerRepo.RemoveLike(answer, user)
	return svc.answerRepo.GetLikesCount(answer)
}

func (svc *userUsecase) Answer(user *models.User, question *models.Question, answer string) *models.Answer {
	a := user.Answer(question, answer)
	return svc.answerRepo.Save(a)
}
