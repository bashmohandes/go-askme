package models

import "time"

type entity struct {
	ID        string
	CreatedOn time.Time
}

//UserID type
type UserID string

//Question asked by users
type Question struct {
	entity
	Text    string
	AskedBy UserID
	AskedTo UserID
}

//Answer by me
type Answer struct {
	entity
	Text       string
	Likes      int
	AnsweredBy UserID
	QuestionID string
}

//LoadQuestions loads all questions for a specific user
func LoadQuestions(user UserID) []Question {
	result := []Question{
		Question{
			entity: entity{
				ID:        "2332323",
				CreatedOn: time.Now(),
			},
			AskedBy: "TestUser1",
			AskedTo: "Bashmohandes",
			Text:    "Who Are You?",
		},
		Question{
			entity: entity{
				ID:        "2362323",
				CreatedOn: time.Now(),
			},
			AskedBy: "TestUser2",
			AskedTo: "Bashmohandes",
			Text:    "What Are You?",
		},
	}

	return result
}
