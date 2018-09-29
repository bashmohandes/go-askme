package repository

import "github.com/bashmohandes/go-askme/internal/domain"

// AnswerRepository represents the basic answer repo functionality
type AnswerRepository interface {
	LoadAnswers(userID models.UniqueID) []models.Answer
	AddLike(answer *models.Answer, user *models.User) uint
}
