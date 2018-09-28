package models

import (
	"time"
)

// User type
type User struct {
	Entity
	Email string
	Name  string
}

// Answer the specified question
func (user *User) Answer(q *Question, text string) *Answer {
	return &Answer{
		UserEntity: UserEntity{
			Entity: Entity{
				ID:        NewUniqueID(),
				CreatedOn: time.Now(),
			},
			CreatedBy: *user,
		},
		Likes: 0,
		Text:  text,
	}
}
