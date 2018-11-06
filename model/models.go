package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//Answer by me
type Answer struct {
	gorm.Model
	Text       string
	Likes      uint
	QuestionID uint
	Question   Question
	LikedBy    map[uint]bool
	UserID     uint
	User       User
}

//Question asked by users
type Question struct {
	gorm.Model
	ToUser     User
	ToUserID   uint
	Text       string
	AnswerID   *uint
	FromUser   User
	FromUserID uint
}

// User type
type User struct {
	gorm.Model
	Email          string
	Name           string
	HashedPassword []byte
	Answers        []Answer
	Questions      []Question
}

// Answer the specified question
func (user *User) Answer(q *Question, answer string) *Answer {
	return &Answer{
		Likes:      0,
		Text:       answer,
		QuestionID: q.ID,
		User:       *user,
		UserID:     user.ID,
		LikedBy:    make(map[uint]bool),
	}
}

// Ask creates a new question that is asked from user to user askee
func (user *User) Ask(other *User, question string) *Question {
	return &Question{
		Text:       question,
		ToUser:     *other,
		ToUserID:   user.ID,
		FromUser:   *user,
		FromUserID: user.ID,
	}
}

// Verify user password
func (user *User) Verify(password string) bool {
	return bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)) == nil
}

// NewUser creates a new user
func NewUser(email string, name string, password string) (*User, error) {
	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:          email,
		Name:           name,
		HashedPassword: hpass,
	}, nil
}
