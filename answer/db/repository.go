package answer

import (
	"github.com/bashmohandes/go-askme/answer"
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/web/framework"
)

type answersRepo struct {
	framework.Connection
}

// LoadAnswers loads the specified user's set of answers
func (r *answersRepo) LoadAnswers(userID uint) ([]*models.Answer, error) {
	db, err := r.Connect()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	var answers []*models.Answer
	err = db.Preload("User").Preload("Question").Find(&answers).Error
	if err != nil {
		return nil, err
	}
	return answers, nil
}

// AddLike adds a like to the specified answer
func (r *answersRepo) AddLike(answer *models.Answer, user *models.User) {

}

func (r *answersRepo) RemoveLike(answer *models.Answer, user *models.User) {

}

func (r *answersRepo) GetLikesCount(answer *models.Answer) uint {
	return 0
}

func (r *answersRepo) Add(answer *models.Answer) (*models.Answer, error) {
	db, err := r.Connect()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	err = db.Create(answer).Error
	if err != nil {
		return nil, err
	}

	return answer, nil
}

// NewRepository creates a new repo object
func NewRepository(conn framework.Connection) answer.Repository {
	return &answersRepo{conn}
}
