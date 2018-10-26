package user

import "github.com/bashmohandes/go-askme/model"

// Repository defines the basic repo operations
type Repository interface {
	Add(user *models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(id models.UniqueID) (*models.User, error)
}
