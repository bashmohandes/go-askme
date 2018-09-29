package repository

import (
	"time"

	"github.com/bashmohandes/go-askme/internal/repository"

	"github.com/bashmohandes/go-askme/internal/domain"
)

// InMemoryQuestionsRepo is a test repo that saves data in memory
type questionsRepo struct {
	data map[models.UniqueID]*models.Question
}

// LoadQuestions loads the specified user's set of questions
func (r *questionsRepo) LoadQuestions(userID models.UniqueID) []*models.Question {
	result := []*models.Question{
		&models.Question{
			UserEntity: models.UserEntity{
				Entity: models.Entity{
					CreatedOn: time.Now(),
					ID:        models.NewUniqueID(),
				},
				CreatedBy: &models.User{
					Entity: models.Entity{
						CreatedOn: time.Now(),
						ID:        models.NewUniqueID(),
					},
					Email: "Bashmohandes@live.com",
					Name:  "Mohamed Elsherif",
				},
			},
			Text: "Who Are You?",
		},
		&models.Question{
			UserEntity: models.UserEntity{
				Entity: models.Entity{
					CreatedOn: time.Now(),
					ID:        models.NewUniqueID(),
				},
				CreatedBy: &models.User{
					Entity: models.Entity{
						CreatedOn: time.Now(),
						ID:        models.NewUniqueID(),
					},
					Email: "Bashmohandes@live.com",
					Name:  "Mohamed Elsherif",
				},
			},
			Text: "Where Are You?",
		},
	}

	return result
}

// Save the question specified
func (r *questionsRepo) Save(q *models.Question) *models.Question {
	r.data[q.ID] = q
	return q
}

// NewQuestionRepository creates a new repo object
func NewQuestionRepository() repository.QuestionRepository {
	return &questionsRepo{
		data: make(map[models.UniqueID]*models.Question),
	}
}
