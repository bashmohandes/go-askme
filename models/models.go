package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//Answer by me
type Answer struct {
	gorm.Model
	Text       string
	QuestionID uint `gorm:"type:int REFERENCES questions(id)"`
	Question   Question
	UserID     uint `gorm:"type:int REFERENCES users(id)"`
	User       User
}

//Question asked by users
type Question struct {
	gorm.Model
	ToUser     User
	ToUserID   uint `gorm:"type:int REFERENCES users(id)"`
	Text       string
	AnswerID   *uint `gorm:"default: null"`
	FromUser   User
	FromUserID uint `gorm:"type:int REFERENCES users(id)"`
}

// User type
type User struct {
	gorm.Model
	Email             string `gorm:"type:varchar(100);unique_index"`
	Name              string
	HashedPassword    []byte
	Answers           []Answer
	QuestionsReceived []Question `gorm:"FOREIGNKEY:ToUserID"`
	QuestionsSent     []Question `gorm:"FOREIGNKEY:FromUserID"`
	FollowedUsers     []*User    `gorm:"many2many:followers;association_jointable_foreignkey:follower_user_id"`
}

// Answer the specified question
func (user *User) Answer(q *Question, answer string) *Answer {
	return &Answer{
		Text:       answer,
		QuestionID: q.ID,
		User:       *user,
		UserID:     user.ID,
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

// Follow the specified user
func (user *User) Follow(other *User) {
	user.FollowedUsers = append(user.FollowedUsers, other)
}

// Unfollow the specified user
func (user *User) Unfollow(other *User) {

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
