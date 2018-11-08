package question

import (
	"github.com/bashmohandes/go-askme/model"
	"github.com/bashmohandes/go-askme/question"
	"github.com/bashmohandes/go-askme/web/framework"
)

type questionsRepo struct {
	framework.Connection
}

// LoadUnansweredQuestions loads the specified user's set of questions
func (r *questionsRepo) LoadUnansweredQuestions(userID uint) ([]*models.Question, error) {
	db, err := r.Connect()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var result []*models.Question
	err = db.Where("to_user_id = ? AND answer_id is NULL", userID).Preload("FromUser").Order("created_at desc").Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Save the question specified
func (r *questionsRepo) Add(q *models.Question) (*models.Question, error) {
	db, err := r.Connect()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	if q.ID > 0 {
		err = db.Save(q).Error
	} else {
		err = db.Create(q).Error
	}

	if err != nil {
		return nil, err
	}

	return q, nil
}

func (r *questionsRepo) GetByID(id uint) (*models.Question, error) {
	db, err := r.Connect()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var question models.Question
	err = db.Find(&question, id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

// NewRepository creates a new repo object
func NewRepository(conn framework.Connection) question.Repository {
	return &questionsRepo{conn}
}
