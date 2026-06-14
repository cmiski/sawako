package authentication

import (
	"context"
	"time"

	"github.com/cmiski/sawako/services/auth/internal/credential"
	"github.com/cmiski/sawako/services/auth/internal/refreshtoken"
	"github.com/cmiski/sawako/services/auth/internal/user"
	"github.com/google/uuid"
)

type UserService interface {
	Create(
		ctx context.Context,
		user *user.User,
	) error
	GetByEmail(
		ctx context.Context,
		email string,
	) (*user.User, error)
}

type CredentialRepository interface {
	Create(
		ctx context.Context,
		credential *credential.Credential,
	) error
	GetByUserIDAndType(
		ctx context.Context,
		userID uuid.UUID,
		credentialType credential.CredentialType,
	) (*credential.Credential, error)
}

type RefreshTokenRepository interface {
	Create(
		ctx context.Context,
		token *refreshtoken.RefreshToken,
	) error

	GetByHash(
		ctx context.Context,
		tokenHash string,
	) (*refreshtoken.RefreshToken, error)

	Revoke(
		ctx context.Context,
		id uuid.UUID,
	) error

	RevokeByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) error
}

type RefreshTokenService interface {
	New(
		userID uuid.UUID,
		tokenHash string,
		expiresAt time.Time,
	) *refreshtoken.RefreshToken
}
