package models

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("test@test.com", "Test User", "p@ssw0rd")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if user.ID != 0 {
		t.Error("ID should be 0")
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

	if user.HashedPassword == nil || len(user.HashedPassword) == 0 {
		t.Error("Password is nor hashed")
		t.Fail()
	}
}

func TestAskUserQuestion(t *testing.T) {
	user1, _ := NewUser("test1@test.com", "Test User1", "p@ssws0rd")
	user2, _ := NewUser("test2@test.com", "Test User2", "p@ssw0rd1")
	question := user1.Ask(user2, "This is a dumb question!")

	if question == nil {
		t.Error("Nil question")
		t.Fail()
	}

	if question.ID != 0 {
		t.Error("Question ID should be 0")
		t.Fail()
	}

	if question.FromUser.ID != user1.ID {
		t.Errorf("Wrong asker %d!= %d", question.FromUser.ID, user1.ID)
		t.Fail()
	}

	if question.ToUser.ID != user2.ID {
		t.Errorf("Wrong askee %d != %d", question.ToUser.ID, user2.ID)
		t.Fail()
	}

	if question.Text != "This is a dumb question!" {
		t.Error("Missing question text")
		t.Fail()
	}
}

func TestAnswerQuestion(t *testing.T) {
	user1, _ := NewUser("test1@test.com", "Test User1", "p@ssw0rd")
	user2, _ := NewUser("test2@test.com", "Test User2", "p@ssw0rd")
	question := user1.Ask(user2, "This is a dumb question!")

	answer := user2.Answer(question, "This is a snarky answer!")

	if answer == nil {
		t.Error("Nil answer")
		t.Fail()
	}

	if answer.ID != 0 {
		t.Error("ID should be 0")
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
