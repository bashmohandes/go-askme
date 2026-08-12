package user

import (
	"errors"
	"testing"

	"github.com/bashmohandes/go-askme/answer/inmemory"
	"github.com/bashmohandes/go-askme/models"
	"github.com/bashmohandes/go-askme/question/inmemory"
	user "github.com/bashmohandes/go-askme/user/inmemory"
	"github.com/jinzhu/gorm"
)

func TestAsk(t *testing.T) {
	sut := NewAsksUsecase(question.NewRepository(), answer.NewRepository(), user.NewRepository())
	from := &models.User{}
	from.ID = 123
	to := &models.User{}
	to.ID = 45
	question := sut.Ask(from, to, "test question")
	if question.FromUser.ID != from.ID ||
		question.ToUser.ID != to.ID ||
		question.Text != "test question" {
		t.Fail()
	}
}

func TestAnswerRejectsUserWhoIsNotQuestionRecipient(t *testing.T) {
	sut := NewAnswersUsecase(question.NewRepository(), answer.NewRepository(), user.NewRepository())
	question := &models.Question{ToUserID: 42}

	got, err := sut.Answer(&models.User{Model: gorm.Model{ID: 7}}, question, "not allowed")

	if got != nil {
		t.Errorf("Answer() = %#v, want nil", got)
	}
	if !errors.Is(err, ErrNotQuestionRecipient) {
		t.Errorf("Answer() error = %v, want %v", err, ErrNotQuestionRecipient)
	}
}

func TestAnswerSavesAnswerForQuestionRecipient(t *testing.T) {
	questions := question.NewRepository()
	sut := NewAnswersUsecase(questions, answer.NewRepository(), user.NewRepository())
	recipient := &models.User{Model: gorm.Model{ID: 42}}
	question := &models.Question{
		Model:    gorm.Model{ID: 10},
		ToUser:   *recipient,
		ToUserID: recipient.ID,
	}

	got, err := sut.Answer(recipient, question, "an answer")

	if err != nil {
		t.Fatalf("Answer() error = %v", err)
	}
	if got.Text != "an answer" || got.UserID != recipient.ID || got.QuestionID != question.ID {
		t.Errorf("Answer() = %#v, want recipient and question fields", got)
	}
	if question.AnswerID == nil {
		t.Fatal("Answer() did not mark question as answered")
	}
}
