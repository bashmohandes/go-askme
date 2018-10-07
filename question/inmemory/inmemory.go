package question

import (
	"time"

	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/question"
)

// InMemoryQuestionsRepo is a test repo that saves data in memory
type questionsRepo struct {
	data map[models.UniqueID]*models.Question
}

// LoadUnansweredQuestions loads the specified user's set of questions
func (r *questionsRepo) LoadUnansweredQuestions(userID models.UniqueID) []*models.Question {
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

// NewRepository creates a new repo object
func NewRepository() question.Repository {
	return &questionsRepo{
		data: make(map[models.UniqueID]*models.Question),
	}
}
