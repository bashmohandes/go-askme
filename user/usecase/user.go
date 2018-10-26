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
	FindUserByEmail(email string) *models.User
}

// AnswersUsecase for the registered user
type AnswersUsecase interface {
	FetchUnansweredQuestions(userID string) (*QuestionsFeed, error)
	Answer(user *models.User, question *models.Question, answer string) *models.Answer
}

// AuthUsecase defines authentication use cases
type AuthUsecase interface {
	Signin(email string, password string) (*models.User, error)
	Signup(email string, password string, name string) (*models.User, error)
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
func NewAnswersUsecase(qRepo question.Repository, aRepo answer.Repository) AnswersUsecase {
	return &userUsecase{
		questionRepo: qRepo,
		answerRepo:   aRepo,
	}
}

// NewAuthUsecase creates a new auth usecase
func NewAuthUsecase(uRepo user.Repository) AuthUsecase {
	return &authUsecase{
		userRepo: uRepo,
	}
}

// LoadQuestions load questions model
func (svc *userUsecase) FetchUnansweredQuestions(userID string) (*QuestionsFeed, error) {
	uqID, err := models.ParseUniqueID(userID)
	if err != nil {
		return nil, err
	}
	feed := QuestionsFeed{
		Items: make([]*QuestionFeedItem, 0),
	}
	questions := svc.questionRepo.LoadUnansweredQuestions(uqID)
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

func (svc *userUsecase) FindUserByEmail(email string) *models.User {
	u, _ := svc.userRepo.GetByEmail(email)
	return u
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
