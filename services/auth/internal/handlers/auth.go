package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/cmiski/sawako/services/auth/internal/authentication"
	"github.com/cmiski/sawako/services/auth/internal/user"
)

const minPasswordLength = 8

type AuthHandler struct {
	auth *authentication.Service
}

func NewAuthHandler(
	auth *authentication.Service,
) *AuthHandler {
	return &AuthHandler{
		auth: auth,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if err := validateCredentials(
		req.Email,
		req.Password,
	); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	err := h.auth.Register(
		r.Context(),
		authentication.RegisterRequest{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmailAlreadyExists):
			writeError(
				w,
				http.StatusConflict,
				"email already exists",
			)

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"registration failed",
			)
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if err := validateCredentials(
		req.Email,
		req.Password,
	); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	resp, err := h.auth.Login(
		r.Context(),
		authentication.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		if errors.Is(err, authentication.ErrInvalidCredentials) {
			writeError(
				w,
				http.StatusUnauthorized,
				"invalid credentials",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"login failed",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		loginResponse{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
		},
	)
}

func (h *AuthHandler) Refresh(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req refreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"refresh_token is required",
		)
		return
	}

	resp, err := h.auth.Refresh(
		r.Context(),
		authentication.RefreshRequest{
			RefreshToken: req.RefreshToken,
		},
	)
	if err != nil {
		if errors.Is(err, authentication.ErrInvalidRefreshToken) {
			writeError(
				w,
				http.StatusUnauthorized,
				"invalid refresh token",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"refresh failed",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		loginResponse{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
		},
	)
}

func (h *AuthHandler) Logout(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req refreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"refresh_token is required",
		)
		return
	}

	err := h.auth.Logout(
		r.Context(),
		authentication.RefreshRequest{
			RefreshToken: req.RefreshToken,
		},
	)
	if err != nil {
		if errors.Is(err, authentication.ErrInvalidRefreshToken) {
			writeError(
				w,
				http.StatusUnauthorized,
				"invalid refresh token",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"logout failed",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateCredentials(
	email string,
	password string,
) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email is required")
	}

	if len(password) < minPasswordLength {
		return errors.New(
			"password must be at least 8 characters",
		)
	}

	return nil
}
