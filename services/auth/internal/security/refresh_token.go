package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/cmiski/sawako/services/auth/internal/authentication"
)

const refreshTokenBytes = 32

type RandomRefreshTokenGenerator struct{}

func NewRandomRefreshTokenGenerator() authentication.RefreshTokenGenerator {
	return &RandomRefreshTokenGenerator{}
}

func (g *RandomRefreshTokenGenerator) Generate() (string, error) {
	buf := make([]byte, refreshTokenBytes)

	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf(
			"random refresh token: %w",
			err,
		)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
