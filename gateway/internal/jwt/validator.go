package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const claimSubject = "sub"

type Validator struct {
	secret []byte
}

func NewValidator(
	secret string,
) *Validator {
	return &Validator{
		secret: []byte(secret),
	}
}

func (v *Validator) Validate(
	tokenString string,
) (uuid.UUID, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf(
					"jwt: unexpected signing method: %v",
					token.Header["alg"],
				)
			}

			return v.secret, nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"jwt: parse token: %w",
			ErrInvalidToken,
		)
	}

	if !token.Valid {
		return uuid.Nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}

	subject, ok := claims[claimSubject].(string)
	if !ok || subject == "" {
		return uuid.Nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	return userID, nil
}
