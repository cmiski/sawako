package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/cmiski/sawako/services/auth/internal/credential"
	"github.com/cmiski/sawako/services/auth/internal/user"
)

type RegisterRequest struct {
	Email    string
	Password string
}

type Service struct {
	users       UserService
	credentials CredentialRepository
	hasher      PasswordHasher
	txManager   TransactionManager
}

func NewService(
	users UserService,
	credentials CredentialRepository,
	hasher PasswordHasher,
	txManager TransactionManager,
) *Service {
	return &Service{
		users:       users,
		credentials: credentials,
		hasher:      hasher,
		txManager:   txManager,
	}
}

func (s *Service) Register(
	ctx context.Context,
	req RegisterRequest,
) error {
	email := strings.ToLower(
		strings.TrimSpace(req.Email),
	)

	passwordHash, err := s.hasher.Hash(
		req.Password,
	)
	if err != nil {
		return fmt.Errorf(
			"register user: hash password: %w",
			err,
		)
	}

	err = s.txManager.WithinTransaction(
		ctx,
		func(ctx context.Context) error {
			u := &user.User{
				Email: email,
			}

			if err := s.users.Create(
				ctx,
				u,
			); err != nil {
				return fmt.Errorf(
					"create user: %w",
					err,
				)
			}

			c := &credential.Credential{
				UserID:         u.ID,
				CredentialType: "password",
				PasswordHash:   passwordHash,
			}

			if err := s.credentials.Create(
				ctx,
				c,
			); err != nil {
				return fmt.Errorf(
					"create credential: %w",
					err,
				)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf(
			"register user: %w",
			err,
		)
	}

	return nil
}
