package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Email           string
	IsEmailVerified bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
