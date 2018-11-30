package user

import "github.com/bashmohandes/go-askme/models"

// Repository defines the basic repo operations
type Repository interface {
	Add(user *models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	Persist(user *models.User) error
}
