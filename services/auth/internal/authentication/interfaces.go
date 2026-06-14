package auth

import (
	"context"

	"github.com/cmiski/sawako/services/auth/internal/credential"
	"github.com/cmiski/sawako/services/auth/internal/user"
)

type UserService interface {
	Create(
		ctx context.Context,
		user *user.User,
	) error
}

type CredentialRepository interface {
	Create(
		ctx context.Context,
		credential *credential.Credential,
	) error
}
