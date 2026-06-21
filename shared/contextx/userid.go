package contextx

import (
	"context"

	"github.com/google/uuid"
)

const UserIDHeader = "X-User-ID"

type userIDKey struct{}

func SetUserID(
	ctx context.Context,
	userID uuid.UUID,
) context.Context {
	return context.WithValue(
		ctx,
		userIDKey{},
		userID,
	)
}

func GetUserID(
	ctx context.Context,
) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return uuid.Nil, false
	}

	return userID, true
}
