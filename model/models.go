package models

import (
	"time"

	"github.com/google/uuid"
)

// UniqueID type
type UniqueID uuid.UUID

// EmptyUniqueID represents empty UniqueID
var EmptyUniqueID = UniqueID(uuid.Nil)

// Entity base
type Entity struct {
	ID        UniqueID
	CreatedOn time.Time
}

// UserEntity base
type UserEntity struct {
	Entity
	CreatedBy *User
}

//Answer by me
type Answer struct {
	UserEntity
	Text       string
	Likes      uint
	QuestionID *UniqueID
	LikedBy    map[UniqueID]bool
}

//Question asked by users
type Question struct {
	UserEntity
	To       *User
	Text     string
	AnswerID *UniqueID
}

// User type
type User struct {
	Entity
	Email string
	Name  string
}

// Answer the specified question
func (user *User) Answer(q *Question, answer string) *Answer {
	return &Answer{
		UserEntity: UserEntity{
			Entity: Entity{
				ID:        NewUniqueID(),
				CreatedOn: time.Now(),
			},
			CreatedBy: user,
		},
		Likes:      0,
		Text:       answer,
		QuestionID: &q.ID,
		LikedBy:    make(map[UniqueID]bool),
	}
}

// Ask creates a new question that is asked from user to user askee
func (user *User) Ask(other *User, question string) *Question {
	return &Question{
		UserEntity: UserEntity{
			Entity: Entity{
				ID:        NewUniqueID(),
				CreatedOn: time.Now(),
			},
			CreatedBy: user,
		},
		Text: question,
		To:   other,
	}
}

// NewUser creates a new user
func NewUser(email string, name string) *User {
	return &User{
		Entity: Entity{
			ID:        NewUniqueID(),
			CreatedOn: time.Now(),
		},
		Email: email,
		Name:  name,
	}
}

// NewUniqueID generates new UniqueID
func NewUniqueID() UniqueID {
	return UniqueID(uuid.New())
}
