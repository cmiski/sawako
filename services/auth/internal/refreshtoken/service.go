package refreshtoken

import (
	"time"

	"github.com/cmiski/sawako/shared/uuidx"
	"github.com/google/uuid"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) New(
	userID uuid.UUID,
	tokenHash string,
	expiresAt time.Time,
) *RefreshToken {
	return &RefreshToken{
		ID:        uuidx.NewV7(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
}
