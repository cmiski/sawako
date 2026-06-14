package refreshtoken

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		token *RefreshToken,
	) error

	GetByHash(
		ctx context.Context,
		tokenHash string,
	) (*RefreshToken, error)

	Revoke(
		ctx context.Context,
		id uuid.UUID,
	) error

	RevokeByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) error

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*RefreshToken, error)
}
