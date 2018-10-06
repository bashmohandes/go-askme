package models

import (
	"time"

	"github.com/google/uuid"
)

// UniqueID type
type UniqueID uuid.UUID

// EmptyUniqueID represents empty UniqueID
var EmptyUniqueID = UniqueID(uuid.Nil)

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

func (u UniqueID) String() string {
	return uuid.UUID(u).String()
}
