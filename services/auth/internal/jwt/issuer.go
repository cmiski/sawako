package jwt

import (
	"fmt"
	"time"

	"github.com/cmiski/sawako/services/auth/internal/authentication"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

const (
	claimSubject = "sub"
)

type Issuer struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewIssuer(
	secret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) authentication.TokenIssuer {
	return &Issuer{
		secret:          []byte(secret),
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (i *Issuer) IssueAccessToken(
	claims authentication.Claims,
) (string, error) {
	return i.issueToken(
		claims,
		i.accessTokenTTL,
	)
}

func (i *Issuer) IssueRefreshToken(
	claims authentication.Claims,
) (string, error) {
	return i.issueToken(
		claims,
		i.refreshTokenTTL,
	)
}

func (i *Issuer) issueToken(
	claims authentication.Claims,
	ttl time.Duration,
) (string, error) {
	now := time.Now()

	token := jwtlib.NewWithClaims(
		jwtlib.SigningMethodHS256,
		jwtlib.MapClaims{
			claimSubject: claims.UserID.String(),
			"iat":        now.Unix(),
			"exp":        now.Add(ttl).Unix(),
		},
	)

	signed, err := token.SignedString(i.secret)
	if err != nil {
		return "", fmt.Errorf(
			"jwt: sign token: %w",
			err,
		)
	}

	return signed, nil
}
