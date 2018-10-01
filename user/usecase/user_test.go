package usecase

import (
	"testing"

	answerRepo "github.com/bashmohandes/go-askme/answer/repository"
	"github.com/bashmohandes/go-askme/model"
	questionRepo "github.com/bashmohandes/go-askme/question/repository"
)

func TestAsk(t *testing.T) {
	sut := NewUsecase(questionRepo.NewRepository(), answerRepo.NewRepository())
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
	sut := NewUsecase(questionRepo.NewRepository(), answerRepo.NewRepository())
	user1 := &models.User{}
	user1.ID = models.NewUniqueID()
	user2 := &models.User{}
	user2.ID = models.NewUniqueID()
	to := &models.User{}
	to.ID = models.NewUniqueID()
	question := sut.Ask(user1, to, "test question")
	answer := sut.Answer(user1, question, "test answer")
	count := sut.Like(user1, answer)

	if count != 1 {
		t.Fail()
		t.Errorf("Expected 1, Actual %v", count)
	}

	count = sut.Like(user2, answer)
	if count != 2 {
		t.Fail()
		t.Errorf("Expected 2, Actual %v", count)
	}

	count = sut.Like(user2, answer)
	if count != 2 {
		t.Fail()
		t.Errorf("Expected 2, Actual %v", count)
	}
}

func TestUnlike(t *testing.T) {
	sut := NewUsecase(questionRepo.NewRepository(), answerRepo.NewRepository())
	user1 := &models.User{}
	user1.ID = models.NewUniqueID()
	user2 := &models.User{}
	user2.ID = models.NewUniqueID()
	to := &models.User{}
	to.ID = models.NewUniqueID()
	question := sut.Ask(user1, to, "test question")
	answer := sut.Answer(user1, question, "test answer")
	count := sut.Like(user1, answer)

	if count != 1 {
		t.Fail()
		t.Errorf("Expected 1, Actual %v", count)
	}

	count = sut.Like(user2, answer)
	if count != 2 {
		t.Fail()
		t.Errorf("Expected 2, Actual %v", count)
	}

	count = sut.Unlike(user1, answer)
	if count != 1 {
		t.Fail()
		t.Errorf("Expected 1, Actual %v", count)
	}
}
