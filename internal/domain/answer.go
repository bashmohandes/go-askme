package models

//Answer by me
type Answer struct {
	UserEntity
	Text       string
	Likes      uint
	QuestionID UniqueID
}

// AnswerRepository represents the basic answer repo functionality
type AnswerRepository interface {
	LoadAnswers(userID UniqueID) []Answer
}
