package authentication

import "github.com/google/uuid"

type Claims struct {
	UserID uuid.UUID
}

type TokenIssuer interface {
	IssueAccessToken(
		claims Claims,
	) (string, error)

	IssueRefreshToken(
		claims Claims,
	) (string, error)
}
