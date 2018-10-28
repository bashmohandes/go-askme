package question

import (
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/question"
)

// InMemoryQuestionsRepo is a test repo that saves data in memory
type questionsRepo struct {
	data         map[models.UniqueID]*models.Question
	userQuestion map[models.UniqueID][]*models.Question
}

// LoadUnansweredQuestions loads the specified user's set of questions
func (r *questionsRepo) LoadUnansweredQuestions(userID models.UniqueID) []*models.Question {
	result := make([]*models.Question, 0, len(r.data))
	for _, uq := range r.userQuestion[userID] {
		if uq.AnswerID == nil {
			result = append(result, uq)
		}
	}
	return result
}

// Save the question specified
func (r *questionsRepo) Add(q *models.Question) {
	r.data[q.ID] = q
	if r.userQuestion[q.To.ID] == nil {
		r.userQuestion[q.To.ID] = make([]*models.Question, 0)
	}
	r.userQuestion[q.To.ID] = append(r.userQuestion[q.To.ID], q)
}

// NewRepository creates a new repo object
func NewRepository() question.Repository {
	return &questionsRepo{
		data:         make(map[models.UniqueID]*models.Question),
		userQuestion: make(map[models.UniqueID][]*models.Question),
	}
}
