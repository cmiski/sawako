package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/cmiski/sawako/shared/uuidx"
)

type Service struct {
	repo Repository
}

func NewService(
	repo Repository,
) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	user *User,
) error {
	user.ID = uuidx.NewV7()

	_, err := s.repo.GetByEmail(
		ctx,
		user.Email,
	)

	switch {
	case err == nil:
		return fmt.Errorf(
			"create user: %w",
			ErrEmailAlreadyExists,
		)

	case errors.Is(err, ErrUserNotFound):
		// continue

	default:
		return fmt.Errorf(
			"create user: %w",
			err,
		)
	}

	if err := s.repo.Create(
		ctx,
		user,
	); err != nil {
		return fmt.Errorf(
			"create user: %w",
			err,
		)
	}

	return nil
}

func (s *Service) GetByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	user, err := s.repo.GetByEmail(
		ctx,
		email,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get user by email: %w",
			err,
		)
	}

	return user, nil
}
