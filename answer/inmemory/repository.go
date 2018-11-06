package answer

import (
	"github.com/bashmohandes/go-askme/answer"
	"github.com/bashmohandes/go-askme/model"
)

type answersRepo struct {
	data map[uint]*models.Answer
}

// LoadAnswers loads the specified user's set of answers
func (r *answersRepo) LoadAnswers(userID uint) []*models.Answer {
	result := make([]*models.Answer, 0, len(r.data))
	return result
}

// AddLike adds a like to the specified answer
func (r *answersRepo) AddLike(answer *models.Answer, user *models.User) {

}

func (r *answersRepo) RemoveLike(answer *models.Answer, user *models.User) {

}

func (r *answersRepo) GetLikesCount(answer *models.Answer) uint {
	return 0
}

func (r *answersRepo) Add(answer *models.Answer) {
	r.data[answer.ID] = answer
}

// NewRepository creates a new repo object
func NewRepository() answer.Repository {
	return &answersRepo{
		data: make(map[uint]*models.Answer),
	}
}
