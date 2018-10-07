package user

import (
	"testing"

	"github.com/bashmohandes/go-askme/answer/inmemory"
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/question/inmemory"
)

func TestAsk(t *testing.T) {
	sut := NewAsksUsecase(question.NewRepository(), answer.NewRepository())
	from := &models.User{}
	from.ID = models.NewUniqueID()
	to := &models.User{}
	to.ID = models.NewUniqueID()
	question := sut.Ask(from, to, "test question")
	if question.CreatedBy.ID != from.ID ||
		question.To.ID != to.ID ||
		question.Text != "test question" {
		t.Fail()
	}
}

func TestLike(t *testing.T) {
	qRepo, aRepo := question.NewRepository(), answer.NewRepository()
	userAsks, userAnswers := NewAsksUsecase(qRepo, aRepo), NewAnswersUsecase(qRepo, aRepo)
	user1 := &models.User{}
	user1.ID = models.NewUniqueID()
	user2 := &models.User{}
	user2.ID = models.NewUniqueID()
	to := &models.User{}
	to.ID = models.NewUniqueID()
	question := userAsks.Ask(user1, to, "test question")
	answer := userAnswers.Answer(user1, question, "test answer")
	count := userAsks.Like(user1, answer)

	if count != 1 {
		t.Fail()
		t.Errorf("Expected 1, Actual %v", count)
	}

	count = userAsks.Like(user2, answer)
	if count != 2 {
		t.Fail()
		t.Errorf("Expected 2, Actual %v", count)
	}

	count = userAsks.Like(user2, answer)
	if count != 2 {
		t.Fail()
		t.Errorf("Expected 2, Actual %v", count)
	}
}

func TestUnlike(t *testing.T) {
	qRepo, aRepo := question.NewRepository(), answer.NewRepository()
	userAsks, userAnswers := NewAsksUsecase(qRepo, aRepo), NewAnswersUsecase(qRepo, aRepo)
	user1 := &models.User{}
	user1.ID = models.NewUniqueID()
	user2 := &models.User{}
	user2.ID = models.NewUniqueID()
	to := &models.User{}
	to.ID = models.NewUniqueID()
	question := userAsks.Ask(user1, to, "test question")
	answer := userAnswers.Answer(user1, question, "test answer")
	count := userAsks.Like(user1, answer)

	if count != 1 {
		t.Fail()
		t.Errorf("Expected 1, Actual %v", count)
	}

	count = userAsks.Like(user2, answer)
	if count != 2 {
		t.Fail()
		t.Errorf("Expected 2, Actual %v", count)
	}

	count = userAsks.Unlike(user1, answer)
	if count != 1 {
		t.Fail()
		t.Errorf("Expected 1, Actual %v", count)
	}
}
