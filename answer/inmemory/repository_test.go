package answer

import (
	"testing"

	"github.com/bashmohandes/go-askme/models"
	"github.com/jinzhu/gorm"
)

func TestLoadAnswersFiltersByUser(t *testing.T) {
	repository := NewRepository()
	_, _ = repository.Add(&models.Answer{Model: gorm.Model{ID: 1}, UserID: 10})
	_, _ = repository.Add(&models.Answer{Model: gorm.Model{ID: 2}, UserID: 20})
	_, _ = repository.Add(&models.Answer{Model: gorm.Model{ID: 3}, UserID: 10})

	answers, err := repository.LoadAnswers(10)
	if err != nil {
		t.Fatalf("LoadAnswers() error = %v", err)
	}
	if len(answers) != 2 {
		t.Fatalf("LoadAnswers() returned %d answers, want 2", len(answers))
	}
	for _, answer := range answers {
		if answer.UserID != 10 {
			t.Errorf("LoadAnswers() included answer for user %d", answer.UserID)
		}
	}
}
