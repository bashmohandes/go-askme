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
func (user *User) Answer(q *Question, answer string) *Answer {
	return &Answer{
		UserEntity: UserEntity{
			Entity: Entity{
				ID:        NewUniqueID(),
				CreatedOn: time.Now(),
			},
			CreatedBy: user,
		},
		Likes: 0,
		Text:  answer,
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
	}
}

// Likes records a user liking an answer, and returns the new number of likes
func (user *User) Likes(answer *Answer) uint {
	answer.Likes++
	return answer.Likes
}
