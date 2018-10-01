package models

//Answer by me
type Answer struct {
	UserEntity
	Text       string
	Likes      uint
	QuestionID UniqueID
	LikedBy    map[UniqueID]bool
}
