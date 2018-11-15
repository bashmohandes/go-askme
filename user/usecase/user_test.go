package user

import (
	"testing"

	"github.com/bashmohandes/go-askme/answer/inmemory"
	"github.com/bashmohandes/go-askme/models"
	"github.com/bashmohandes/go-askme/question/inmemory"
	user "github.com/bashmohandes/go-askme/user/inmemory"
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
