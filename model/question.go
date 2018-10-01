package models

//Question asked by users
type Question struct {
	UserEntity
	To       *User
	Text     string
	AnswerID *Answer
}
