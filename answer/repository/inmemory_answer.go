package repository

import (
	"github.com/bashmohandes/go-askme/answer"
	"github.com/bashmohandes/go-askme/model"
)

type answersRepo struct{}

// LoadAnswers loads the specified user's set of answers
func (r *answersRepo) LoadAnswers(userID models.UniqueID) []models.Answer {
	return nil
}

// AddLike adds a like to the specified answer
func (r *answersRepo) AddLike(answer *models.Answer, user *models.User) uint {
	return answer.Likes + 1
}

// NewRepository creates a new repo object
func NewRepository() answer.Repository {
	return &answersRepo{}
}
