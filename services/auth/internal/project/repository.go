package project

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(
		ctx context.Context,
		project *Project,
	) error

	ListByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]Project, error)

	GetByID(
		ctx context.Context,
		id uuid.UUID,
	) (*Project, error)

	Update(
		ctx context.Context,
		project *Project,
	) error

	Delete(
		ctx context.Context,
		id uuid.UUID,
	) error
}
