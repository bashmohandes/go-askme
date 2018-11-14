package framework

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Connection interface
type Connection interface {
	Connect() (*gorm.DB, error)
}

// Repo defines basic repo operation
type repo struct {
	config *Config
}

// NewConnection creates a new connection
func NewConnection(config *Config) Connection {
	return &repo{config}
}

// Connect the app to db
func (r *repo) Connect() (*gorm.DB, error) {
	config := r.config
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", config.PostgresUser, config.PostgresPassword, config.PostgresHost, 5432, config.PostgresDB)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	db.LogMode(config.Debug)
	return db, nil
}
