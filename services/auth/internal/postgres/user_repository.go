package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/cmiski/sawako/services/auth/internal/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(
	pool *pgxpool.Pool,
) user.Repository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	u *user.User,
) error {
	_, err := querier(ctx, r.pool).Exec(
		ctx,
		`
		INSERT INTO users (
			id,
			email,
			is_email_verified
		) VALUES ($1, $2, $3)
		`,
		u.ID,
		u.Email,
		u.IsEmailVerified,
	)
	if err != nil {
		return fmt.Errorf(
			"postgres: create user: %w",
			err,
		)
	}

	return nil
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*user.User, error) {
	row := querier(ctx, r.pool).QueryRow(
		ctx,
		`
		SELECT
			id,
			email,
			is_email_verified,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
		`,
		email,
	)

	var u user.User

	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.IsEmailVerified,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	switch {
	case err == nil:
		return &u, nil

	case errors.Is(err, pgx.ErrNoRows):
		return nil, user.ErrUserNotFound

	default:
		return nil, fmt.Errorf(
			"postgres: get user by email: %w",
			err,
		)
	}
}
