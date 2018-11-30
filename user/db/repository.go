package user

import (
	"github.com/bashmohandes/go-askme/models"
	"github.com/bashmohandes/go-askme/user"
	"github.com/bashmohandes/go-askme/web/framework"
)

type userRepo struct {
	framework.Connection
}

// NewRepository creates a new user repo
func NewRepository(conn framework.Connection) user.Repository {
	return &userRepo{conn}
}

func (u *userRepo) Add(user *models.User) (*models.User, error) {
	db, err := u.Connect()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	err = db.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepo) GetByEmail(email string) (*models.User, error) {
	db, err := u.Connect()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var user models.User
	err = db.Where(&models.User{Email: email}).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) GetByID(id uint) (*models.User, error) {
	db, err := u.Connect()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var user models.User
	err = db.Find(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepo) Persist(user *models.User) error {
	db, err := u.Connect()
	if err != nil {
		return err
	}
	err = db.Save(user).Error
	if err != nil {
		return err
	}

	return nil
}
