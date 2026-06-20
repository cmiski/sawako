package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/cmiski/sawako/services/auth/internal/credential"
	"github.com/cmiski/sawako/shared/uuidx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CredentialRepository struct {
	pool *pgxpool.Pool
}

func NewCredentialRepository(
	pool *pgxpool.Pool,
) credential.Repository {
	return &CredentialRepository{
		pool: pool,
	}
}

func (r *CredentialRepository) Create(
	ctx context.Context,
	c *credential.Credential,
) error {
	if c.ID == uuid.Nil {
		c.ID = uuidx.NewV7()
	}

	_, err := querier(ctx, r.pool).Exec(
		ctx,
		`
		INSERT INTO credentials (
			id,
			user_id,
			credential_type,
			password_hash
		) VALUES ($1, $2, $3, $4)
		`,
		c.ID,
		c.UserID,
		string(c.CredentialType),
		c.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf(
			"postgres: create credential: %w",
			err,
		)
	}

	return nil
}

func (r *CredentialRepository) GetByUserIDAndType(
	ctx context.Context,
	userID uuid.UUID,
	credentialType credential.CredentialType,
) (*credential.Credential, error) {
	row := querier(ctx, r.pool).QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			credential_type,
			password_hash,
			created_at,
			updated_at
		FROM credentials
		WHERE user_id = $1
		  AND credential_type = $2
		`,
		userID,
		string(credentialType),
	)

	var c credential.Credential

	var credentialTypeValue string

	err := row.Scan(
		&c.ID,
		&c.UserID,
		&credentialTypeValue,
		&c.PasswordHash,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	switch {
	case err == nil:
		c.CredentialType = credential.CredentialType(
			credentialTypeValue,
		)
		return &c, nil

	case errors.Is(err, pgx.ErrNoRows):
		return nil, credential.ErrCredentialNotFound

	default:
		return nil, fmt.Errorf(
			"postgres: get credential: %w",
			err,
		)
	}
}
