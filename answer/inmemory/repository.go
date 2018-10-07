package answer

import (
	"fmt"

	"github.com/bashmohandes/go-askme/answer"
	"github.com/bashmohandes/go-askme/model"
)

type answersRepo struct {
	data map[models.UniqueID]*models.Answer
}

// LoadAnswers loads the specified user's set of answers
func (r *answersRepo) LoadAnswers(userID models.UniqueID) []models.Answer {
	return nil
}

// AddLike adds a like to the specified answer
func (r *answersRepo) AddLike(answer *models.Answer, user *models.User) {
	a, ok := r.data[answer.ID]
	if !ok {
		panic(fmt.Sprintf("Answer with %v doesn't exist.", answer.ID))
	}
	if a.LikedBy == nil {
		a.LikedBy = make(map[models.UniqueID]bool)
	}
	a.LikedBy[user.ID] = true
}

func (r *answersRepo) RemoveLike(answer *models.Answer, user *models.User) {
	a, ok := r.data[answer.ID]
	if !ok {
		panic(fmt.Sprintf("Answer with %v doesn't exist.", answer.ID))
	}
	delete(a.LikedBy, user.ID)
}

func (r *answersRepo) GetLikesCount(answer *models.Answer) uint {
	return uint(len(r.data[answer.ID].LikedBy))
}

func (r *answersRepo) Add(answer *models.Answer) {
	r.data[answer.ID] = answer
}

// NewRepository creates a new repo object
func NewRepository() answer.Repository {
	return &answersRepo{
		data: make(map[models.UniqueID]*models.Answer),
	}
}
