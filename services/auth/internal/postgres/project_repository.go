package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/cmiski/sawako/services/auth/internal/project"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	pool *pgxpool.Pool
}

func NewProjectRepository(
	pool *pgxpool.Pool,
) project.Repository {
	return &ProjectRepository{
		pool: pool,
	}
}

func (r *ProjectRepository) Create(
	ctx context.Context,
	p *project.Project,
) error {
	_, err := querier(ctx, r.pool).Exec(
		ctx,
		`
		INSERT INTO projects (
			id,
			user_id,
			name,
			description
		) VALUES ($1, $2, $3, $4)
		`,
		p.ID,
		p.UserID,
		p.Name,
		p.Description,
	)
	if err != nil {
		return fmt.Errorf(
			"postgres: create project: %w",
			err,
		)
	}

	return nil
}

func (r *ProjectRepository) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]project.Project, error) {
	rows, err := querier(ctx, r.pool).Query(
		ctx,
		`
		SELECT
			id,
			user_id,
			name,
			description,
			created_at,
			updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
		`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"postgres: list projects: %w",
			err,
		)
	}
	defer rows.Close()

	var projects []project.Project

	for rows.Next() {
		var p project.Project

		if err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"postgres: scan project: %w",
				err,
			)
		}

		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"postgres: list projects: %w",
			err,
		)
	}

	if projects == nil {
		projects = []project.Project{}
	}

	return projects, nil
}

func (r *ProjectRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*project.Project, error) {
	row := querier(ctx, r.pool).QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			name,
			description,
			created_at,
			updated_at
		FROM projects
		WHERE id = $1
		`,
		id,
	)

	var p project.Project

	err := row.Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.Description,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	switch {
	case err == nil:
		return &p, nil

	case errors.Is(err, pgx.ErrNoRows):
		return nil, project.ErrProjectNotFound

	default:
		return nil, fmt.Errorf(
			"postgres: get project: %w",
			err,
		)
	}
}

func (r *ProjectRepository) Update(
	ctx context.Context,
	p *project.Project,
) error {
	tag, err := querier(ctx, r.pool).Exec(
		ctx,
		`
		UPDATE projects
		SET name = $2,
		    description = $3,
		    updated_at = NOW()
		WHERE id = $1
		`,
		p.ID,
		p.Name,
		p.Description,
	)
	if err != nil {
		return fmt.Errorf(
			"postgres: update project: %w",
			err,
		)
	}

	if tag.RowsAffected() == 0 {
		return project.ErrProjectNotFound
	}

	return nil
}

func (r *ProjectRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	tag, err := querier(ctx, r.pool).Exec(
		ctx,
		`
		DELETE FROM projects
		WHERE id = $1
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf(
			"postgres: delete project: %w",
			err,
		)
	}

	if tag.RowsAffected() == 0 {
		return project.ErrProjectNotFound
	}

	return nil
}
