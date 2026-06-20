package security

import (
	"fmt"

	"github.com/cmiski/sawako/services/auth/internal/authentication"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

type BcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() authentication.PasswordHasher {
	return &BcryptPasswordHasher{}
}

func (h *BcryptPasswordHasher) Hash(
	password string,
) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcryptCost,
	)
	if err != nil {
		return "", fmt.Errorf(
			"bcrypt: hash password: %w",
			err,
		)
	}

	return string(hash), nil
}

func (h *BcryptPasswordHasher) Verify(
	password string,
	hash string,
) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
	if err != nil {
		return fmt.Errorf(
			"bcrypt: verify password: %w",
			err,
		)
	}

	return nil
}
