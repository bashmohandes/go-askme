package user

import (
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/user"
)

type userRepo struct {
	data map[string]*models.User
}

// NewRepository creates a new user repo
func NewRepository() user.Repository {
	return &userRepo{
		data: make(map[string]*models.User),
	}
}

func (u *userRepo) Add(user *models.User) (*models.User, error) {
	u.data[user.Email] = user
	return user, nil
}

func (u *userRepo) GetByEmail(email string) (*models.User, error) {
	return u.data[email], nil
}
