package user

import (
	"fmt"

	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/user"
)

type userRepo struct {
	data map[string]*models.User
	byID map[models.UniqueID]*models.User
}

// NewRepository creates a new user repo
func NewRepository() user.Repository {
	return &userRepo{
		data: make(map[string]*models.User),
		byID: make(map[models.UniqueID]*models.User),
	}
}

func (u *userRepo) Add(user *models.User) (*models.User, error) {
	u.data[user.Email] = user
	u.byID[user.ID] = user
	return user, nil
}

func (u *userRepo) GetByEmail(email string) (*models.User, error) {
	user, ok := u.data[email]
	if !ok {
		return nil, fmt.Errorf("User doesn't exist")
	}
	return user, nil
}

func (u *userRepo) GetByID(id models.UniqueID) (*models.User, error) {
	user, ok := u.byID[id]
	if !ok {
		return nil, fmt.Errorf("User doesn't exist")
	}
	return user, nil
}
