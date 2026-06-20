package security

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/cmiski/sawako/services/auth/internal/authentication"
)

type SHA256TokenHasher struct{}

func NewSHA256TokenHasher() authentication.TokenHasher {
	return &SHA256TokenHasher{}
}

func (h *SHA256TokenHasher) Hash(
	token string,
) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
