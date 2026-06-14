package credential

import (
	"context"

	"github.com/google/uuid"
)

type CredentialType string

const (
	CredentialTypePassword CredentialType = "password"
	CredentialTypeGoogle   CredentialType = "google"
	CredentialTypeGitHub   CredentialType = "github"
	CredentialTypePasskey  CredentialType = "passkey"
)

type Repository interface {
	Create(
		ctx context.Context,
		credential *Credential,
	) error

	GetByUserIDAndType(
		ctx context.Context,
		userID uuid.UUID,
		credentialType CredentialType,
	) (*Credential, error)
}
