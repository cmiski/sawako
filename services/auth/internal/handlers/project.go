package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/cmiski/sawako/services/auth/internal/project"
	"github.com/cmiski/sawako/shared/contextx"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projects *project.Service
}

func NewProjectHandler(
	projects *project.Service,
) *ProjectHandler {
	return &ProjectHandler{
		projects: projects,
	}
}

type createProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type projectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (h *ProjectHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := contextx.GetUserID(r.Context())
	if !ok {
		WriteUnauthorized(w, "missing user identity")
		return
	}

	var req createProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	created, err := h.projects.Create(
		r.Context(),
		userID,
		project.CreateInput{
			Name:        req.Name,
			Description: req.Description,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, project.ErrNameRequired):
			writeError(
				w,
				http.StatusBadRequest,
				"name is required",
			)

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"failed to create project",
			)
		}

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		toProjectResponse(created),
	)
}

func (h *ProjectHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := contextx.GetUserID(r.Context())
	if !ok {
		WriteUnauthorized(w, "missing user identity")
		return
	}

	projects, err := h.projects.List(
		r.Context(),
		userID,
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"failed to list projects",
		)
		return
	}

	responses := make(
		[]projectResponse,
		0,
		len(projects),
	)

	for _, p := range projects {
		responses = append(
			responses,
			toProjectResponse(&p),
		)
	}

	writeJSON(
		w,
		http.StatusOK,
		responses,
	)
}

func (h *ProjectHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := contextx.GetUserID(r.Context())
	if !ok {
		WriteUnauthorized(w, "missing user identity")
		return
	}

	projectID, err := uuid.Parse(
		chi.URLParam(r, "id"),
	)
	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid project id",
		)
		return
	}

	var req updateProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	updated, err := h.projects.Update(
		r.Context(),
		userID,
		projectID,
		project.UpdateInput{
			Name:        req.Name,
			Description: req.Description,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, project.ErrProjectNotFound):
			writeError(
				w,
				http.StatusNotFound,
				"project not found",
			)

		case errors.Is(err, project.ErrNoFieldsToUpdate):
			writeError(
				w,
				http.StatusBadRequest,
				"no fields to update",
			)

		case errors.Is(err, project.ErrNameRequired):
			writeError(
				w,
				http.StatusBadRequest,
				"name is required",
			)

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"failed to update project",
			)
		}

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		toProjectResponse(updated),
	)
}

func (h *ProjectHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := contextx.GetUserID(r.Context())
	if !ok {
		WriteUnauthorized(w, "missing user identity")
		return
	}

	projectID, err := uuid.Parse(
		chi.URLParam(r, "id"),
	)
	if err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid project id",
		)
		return
	}

	err = h.projects.Delete(
		r.Context(),
		userID,
		projectID,
	)
	if err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"project not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"failed to delete project",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toProjectResponse(
	p *project.Project,
) projectResponse {
	return projectResponse{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
