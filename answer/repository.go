package answer

import "github.com/bashmohandes/go-askme/model"

// Repository represents the basic answer repo functionality
type Repository interface {
	LoadAnswers(userID models.UniqueID) []models.Answer
	AddLike(answer *models.Answer, user *models.User) uint
}
