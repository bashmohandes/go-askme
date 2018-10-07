package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	user := NewUser("test@test.com", "Test User")
	if user.ID == UniqueID(uuid.Nil) {
		t.Error("Null ID")
		t.Fail()
	}

	if user.Name != "Test User" {
		t.Error("Wrong name")
		t.Fail()
	}

	if user.Email != "test@test.com" {
		t.Error("Wrong email")
		t.Fail()
	}
}

func TestAskUserQuestion(t *testing.T) {
	user1 := NewUser("test1@test.com", "Test User1")
	user2 := NewUser("test2@test.com", "Test User2")
	question := user1.Ask(user2, "This is a dumb question!")

	if question == nil {
		t.Error("Nil question")
		t.Fail()
	}

	if question.ID == EmptyUniqueID {
		t.Error("Null Question ID")
		t.Fail()
	}

	if question.CreatedBy.ID != user1.ID {
		t.Errorf("Wrong asker %s!= %s", question.CreatedBy.ID, user1.ID)
		t.Fail()
	}

	if question.To.ID != user2.ID {
		t.Errorf("Wrong askee %s != %s", question.To.ID, user2.ID)
		t.Fail()
	}

	if question.Text != "This is a dumb question!" {
		t.Error("Missing question text")
		t.Fail()
	}
}

func TestAnswerQuestion(t *testing.T) {
	user1 := NewUser("test1@test.com", "Test User1")
	user2 := NewUser("test2@test.com", "Test User2")
	question := user1.Ask(user2, "This is a dumb question!")

	answer := user2.Answer(question, "This is a snarky answer!")

	if answer == nil {
		t.Error("Nil answer")
		t.Fail()
	}

	if answer.ID == EmptyUniqueID {
		t.Error("Empty ID")
		t.Fail()
	}

	if answer.QuestionID != question.ID {
		t.Error("Wrong question id")
		t.Fail()
	}

	if answer.Text != "This is a snarky answer!" {
		t.Error("Wrong answer")
		t.Fail()
	}
}
