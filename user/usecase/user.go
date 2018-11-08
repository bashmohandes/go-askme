package user

import (
	"fmt"
	"log"
	"time"

	"github.com/bashmohandes/go-askme/answer"
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/question"
	"github.com/bashmohandes/go-askme/user"
)

type userUsecase struct {
	questionRepo question.Repository
	answerRepo   answer.Repository
	userRepo     user.Repository
}

type authUsecase struct {
	userRepo user.Repository
}

// AnswersFeed type
type AnswersFeed struct {
	Items []*AnswerFeedItem
}

// AnswerFeedItem type
type AnswerFeedItem struct {
	AnswerID   uint
	Question   string
	Answer     string
	AnsweredAt time.Time
	Likes      uint
	User       string
	Email      string
}

// QuestionsFeed type
type QuestionsFeed struct {
	Items []*QuestionFeedItem
}

// QuestionFeedItem type
type QuestionFeedItem struct {
	QuestionID uint
	Question   string
	AskedAt    time.Time
	UserID     uint
	User       string
}

// AsksUsecase type
type AsksUsecase interface {
	Ask(from *models.User, to *models.User, question string) *models.Question
	Like(user *models.User, answer *models.Answer) uint
	Unlike(user *models.User, answer *models.Answer) uint
	LoadUserFeed(user *models.User) (*AnswersFeed, error)
	FindUserByEmail(email string) (*models.User, error)
}

// AnswersUsecase for the registered user
type AnswersUsecase interface {
	FetchUnansweredQuestions(email string) (*QuestionsFeed, error)
	Answer(user *models.User, question *models.Question, answer string) *models.Answer
}

// AuthUsecase defines authentication use cases
type AuthUsecase interface {
	Signin(email string, password string) (*models.User, error)
	Signup(email string, password string, name string) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
}

// NewAsksUsecase creates a new service
func NewAsksUsecase(qRepo question.Repository, aRepo answer.Repository, uRepo user.Repository) AsksUsecase {
	return &userUsecase{
		questionRepo: qRepo,
		answerRepo:   aRepo,
		userRepo:     uRepo,
	}
}

// NewAnswersUsecase creates a new service
func NewAnswersUsecase(qRepo question.Repository, aRepo answer.Repository, uRepo user.Repository) AnswersUsecase {
	return &userUsecase{
		questionRepo: qRepo,
		answerRepo:   aRepo,
		userRepo:     uRepo,
	}
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(uRepo user.Repository) AuthUsecase {
	return &authUsecase{
		userRepo: uRepo,
	}
}

// LoadQuestions load questions model
func (svc *userUsecase) FetchUnansweredQuestions(email string) (*QuestionsFeed, error) {
	user, err := svc.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	feed := QuestionsFeed{
		Items: make([]*QuestionFeedItem, 0),
	}
	questions, err := svc.questionRepo.LoadUnansweredQuestions(user.ID)
	if err != nil {
		return nil, err
	}
	for _, q := range questions {
		fi := &QuestionFeedItem{
			AskedAt:    q.CreatedAt,
			QuestionID: q.ID,
			Question:   q.Text,
			User:       q.FromUser.Name,
			UserID:     q.FromUserID,
		}
		feed.Items = append(feed.Items, fi)
	}
	return &feed, nil
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

func (svc *userUsecase) FindUserByEmail(email string) (*models.User, error) {
	return svc.userRepo.GetByEmail(email)
}

func (svc *userUsecase) LoadUserFeed(user *models.User) (*AnswersFeed, error) {
	answers, err := svc.answerRepo.LoadAnswers(user.ID)
	if err != nil {
		return nil, err
	}
	feedItems := make([]*AnswerFeedItem, 0, len(answers))
	for _, a := range answers {
		feedItems = append(feedItems, &AnswerFeedItem{
			Question:   a.Question.Text,
			Answer:     a.Text,
			AnsweredAt: a.CreatedAt,
			User:       a.User.Name,
			Email:      a.User.Email,
		})
	}
	return &AnswersFeed{
		Items: feedItems,
	}, nil
}

func (svc *authUsecase) Signin(email string, password string) (*models.User, error) {
	user, err := svc.userRepo.GetByEmail(email)
	if err != nil {
		log.Println(err.Error())
	}
	if user == nil {
		return nil, fmt.Errorf("Login failed for user %s", email)
	}

	if user.Verify(password) {
		return user, nil
	}

	return nil, fmt.Errorf("Login failed for user %s", email)
}

func (svc *authUsecase) Signup(email string, password string, name string) (*models.User, error) {
	user, err := svc.userRepo.GetByEmail(email)
	if err != nil {
		log.Println(err.Error())
	}
	if user != nil {
		return nil, fmt.Errorf("User with the same email already exists")
	}

	user, err = models.NewUser(email, name, password)
	return svc.userRepo.Add(user)
}

func (svc *authUsecase) FindUserByEmail(email string) (*models.User, error) {
	return svc.userRepo.GetByEmail(email)
}
