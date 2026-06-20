package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/cmiski/sawako/services/auth/internal/refreshtoken"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(
	pool *pgxpool.Pool,
) refreshtoken.Repository {
	return &RefreshTokenRepository{
		pool: pool,
	}
}

func (r *RefreshTokenRepository) Create(
	ctx context.Context,
	token *refreshtoken.RefreshToken,
) error {
	_, err := querier(ctx, r.pool).Exec(
		ctx,
		`
		INSERT INTO refresh_tokens (
			id,
			user_id,
			token_hash,
			expires_at,
			revoked_at
		) VALUES ($1, $2, $3, $4, $5)
		`,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.RevokedAt,
	)
	if err != nil {
		return fmt.Errorf(
			"postgres: create refresh token: %w",
			err,
		)
	}

	return nil
}

func (r *RefreshTokenRepository) GetByHash(
	ctx context.Context,
	tokenHash string,
) (*refreshtoken.RefreshToken, error) {
	return r.getOne(
		ctx,
		`
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			revoked_at,
			created_at,
			updated_at
		FROM refresh_tokens
		WHERE token_hash = $1
		`,
		tokenHash,
	)
}

func (r *RefreshTokenRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*refreshtoken.RefreshToken, error) {
	return r.getOne(
		ctx,
		`
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			revoked_at,
			created_at,
			updated_at
		FROM refresh_tokens
		WHERE id = $1
		`,
		id,
	)
}

func (r *RefreshTokenRepository) Revoke(
	ctx context.Context,
	id uuid.UUID,
) error {
	tag, err := querier(ctx, r.pool).Exec(
		ctx,
		`
		UPDATE refresh_tokens
		SET revoked_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
		  AND revoked_at IS NULL
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"postgres: revoke refresh token: %w",
			err,
		)
	}

	if tag.RowsAffected() == 0 {
		return refreshtoken.ErrRefreshTokenNotFound
	}

	return nil
}

func (r *RefreshTokenRepository) RevokeByUserID(
	ctx context.Context,
	userID uuid.UUID,
) error {
	_, err := querier(ctx, r.pool).Exec(
		ctx,
		`
		UPDATE refresh_tokens
		SET revoked_at = NOW(),
		    updated_at = NOW()
		WHERE user_id = $1
		  AND revoked_at IS NULL
		`,
		userID,
	)
	if err != nil {
		return fmt.Errorf(
			"postgres: revoke refresh tokens by user: %w",
			err,
		)
	}

	return nil
}

func (r *RefreshTokenRepository) getOne(
	ctx context.Context,
	query string,
	arg any,
) (*refreshtoken.RefreshToken, error) {
	row := querier(ctx, r.pool).QueryRow(
		ctx,
		query,
		arg,
	)

	var token refreshtoken.RefreshToken

	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
		&token.UpdatedAt,
	)

	switch {
	case err == nil:
		return &token, nil

	case errors.Is(err, pgx.ErrNoRows):
		return nil, refreshtoken.ErrRefreshTokenNotFound

	default:
		return nil, fmt.Errorf(
			"postgres: get refresh token: %w",
			err,
		)
	}
}
