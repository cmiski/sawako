package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/cmiski/sawako/shared/uuidx"
	"github.com/google/uuid"
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

type CreateInput struct {
	Name        string
	Description string
}

type UpdateInput struct {
	Name        *string
	Description *string
}

func (s *Service) Create(
	ctx context.Context,
	userID uuid.UUID,
	input CreateInput,
) (*Project, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf(
			"create project: %w",
			ErrNameRequired,
		)
	}

	project := &Project{
		ID:          uuidx.NewV7(),
		UserID:      userID,
		Name:        name,
		Description: strings.TrimSpace(input.Description),
	}

	if err := s.repo.Create(
		ctx,
		project,
	); err != nil {
		return nil, fmt.Errorf(
			"create project: %w",
			err,
		)
	}

	created, err := s.repo.GetByID(
		ctx,
		project.ID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create project: load created project: %w",
			err,
		)
	}

	return created, nil
}

func (s *Service) List(
	ctx context.Context,
	userID uuid.UUID,
) ([]Project, error) {
	projects, err := s.repo.ListByUserID(
		ctx,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list projects: %w",
			err,
		)
	}

	return projects, nil
}

func (s *Service) Update(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
	input UpdateInput,
) (*Project, error) {
	project, err := s.ownedProject(
		ctx,
		userID,
		projectID,
	)
	if err != nil {
		return nil, err
	}

	if input.Name == nil && input.Description == nil {
		return nil, fmt.Errorf(
			"update project: %w",
			ErrNoFieldsToUpdate,
		)
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, fmt.Errorf(
				"update project: %w",
				ErrNameRequired,
			)
		}

		project.Name = name
	}

	if input.Description != nil {
		project.Description = strings.TrimSpace(
			*input.Description,
		)
	}

	if err := s.repo.Update(
		ctx,
		project,
	); err != nil {
		return nil, fmt.Errorf(
			"update project: %w",
			err,
		)
	}

	updated, err := s.repo.GetByID(
		ctx,
		project.ID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"update project: load updated project: %w",
			err,
		)
	}

	return updated, nil
}

func (s *Service) Delete(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
) error {
	_, err := s.ownedProject(
		ctx,
		userID,
		projectID,
	)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(
		ctx,
		projectID,
	); err != nil {
		return fmt.Errorf(
			"delete project: %w",
			err,
		)
	}

	return nil
}

func (s *Service) ownedProject(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
) (*Project, error) {
	project, err := s.repo.GetByID(
		ctx,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get project: %w",
			err,
		)
	}

	if project.UserID != userID {
		return nil, ErrProjectNotFound
	}

	return project, nil
}
