package user

import (
	"time"

	"github.com/bashmohandes/go-askme/answer"
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/question"
)

type userUsecase struct {
	questionRepo question.Repository
	answerRepo   answer.Repository
}

// AnswersFeed type
type AnswersFeed struct {
	Items []*AnswerFeedItem
}

// AnswerFeedItem type
type AnswerFeedItem struct {
	AnswerID   models.UniqueID
	Question   string
	Answer     string
	AnsweredAt time.Time
	Likes      uint
	User       string
}

// QuestionsFeed type
type QuestionsFeed struct {
	Items []*QuestionFeedItem
}

// QuestionFeedItem type
type QuestionFeedItem struct {
	QuestionID models.UniqueID
	Question   string
	AskedAt    time.Time
	UserID     models.UniqueID
	User       string
}

// AsksUsecase type
type AsksUsecase interface {
	Ask(from *models.User, to *models.User, question string) *models.Question
	Like(user *models.User, answer *models.Answer) uint
	Unlike(user *models.User, answer *models.Answer) uint
	LoadUserFeed(user *models.User) *AnswersFeed
}

// AnswersUsecase for the registered user
type AnswersUsecase interface {
	FetchUnansweredQuestions(userID models.UniqueID) *QuestionsFeed
	Answer(user *models.User, question *models.Question, answer string) *models.Answer
}

// NewAsksUsecase creates a new service
func NewAsksUsecase(qRepo question.Repository, aRepo answer.Repository) AsksUsecase {
	return &userUsecase{
		questionRepo: qRepo,
		answerRepo:   aRepo,
	}
}

// NewAnswersUsecase creates a new service
func NewAnswersUsecase(qRepo question.Repository, aRepo answer.Repository) AnswersUsecase {
	return &userUsecase{
		questionRepo: qRepo,
		answerRepo:   aRepo,
	}
}

// LoadQuestions load questions model
func (svc *userUsecase) FetchUnansweredQuestions(userID models.UniqueID) *QuestionsFeed {
	feed := QuestionsFeed{
		Items: make([]*QuestionFeedItem, 0),
	}
	questions := svc.questionRepo.LoadUnansweredQuestions(userID)
	for _, q := range questions {
		fi := &QuestionFeedItem{
			AskedAt:    q.CreatedOn,
			QuestionID: q.ID,
			Question:   q.Text,
			User:       q.CreatedBy.Name,
			UserID:     q.CreatedBy.ID,
		}
		feed.Items = append(feed.Items, fi)
	}
	return &feed
}

// Saves question
func (svc *userUsecase) Ask(from *models.User, to *models.User, question string) *models.Question {
	q := from.Ask(to, question)
	svc.questionRepo.Add(q)
	return q
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
	svc.answerRepo.Add(a)
	return a
}

func (svc *userUsecase) LoadUserFeed(user *models.User) *AnswersFeed {
	return &AnswersFeed{
		Items: []*AnswerFeedItem{
			&AnswerFeedItem{
				Question:   "What is your name?",
				Answer:     "Mohamed Elsherif daaah!",
				AnsweredAt: time.Now(),
				Likes:      10,
				User:       "Anonymous",
			},
			&AnswerFeedItem{
				Question:   "What is your age?",
				Answer:     "I would never tell, but it is 35",
				AnsweredAt: time.Now(),
				Likes:      1,
				User:       "Sayed",
			},
		},
	}
}
