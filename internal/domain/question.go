package models

// QuestionRepository interface
type QuestionRepository interface {
	LoadQuestions(user UniqueID) []Question
}

//Question asked by users
type Question struct {
	UserEntity
	Text     string
	AnswerID UniqueID
}

// LoadQuestions load questions model
func LoadQuestions(userID UniqueID) []Question {
	return QuestionsRepo.LoadQuestions(userID)
}

// QuestionsRepo is the reference to the configured repository object
var QuestionsRepo QuestionRepository
