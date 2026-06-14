package credential

import (
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	CredentialType CredentialType
	PasswordHash   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
