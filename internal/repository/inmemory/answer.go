package repository

import (
	"github.com/bashmohandes/go-askme/internal/domain"
	"github.com/bashmohandes/go-askme/internal/repository"
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

// NewAnswerRepository creates a new repo object
func NewAnswerRepository() repository.AnswerRepository {
	return &answersRepo{}
}
