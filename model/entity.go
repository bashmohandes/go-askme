package models

import (
	"time"

	"github.com/google/uuid"
)

// UniqueID type
type UniqueID uuid.UUID

// Entity base
type Entity struct {
	ID        UniqueID
	CreatedOn time.Time
}

// UserEntity base
type UserEntity struct {
	Entity
	CreatedBy *User
}

// NewUniqueID generates new UniqueID
func NewUniqueID() UniqueID {
	return UniqueID(uuid.New())
}
