package repository

import "github.com/bashmohandes/go-askme/internal/domain"

type answersRepo struct{}

// LoadAnswers loads the specified user's set of answers
func (r *answersRepo) LoadAnswers(userID models.UniqueID) []models.Answer {
	return nil
}

// NewAnswerRepository creates a new repo object
func NewAnswerRepository() models.AnswerRepository {
	return &answersRepo{}
}
