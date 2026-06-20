package handlers

import (
	"net/http"
)

type HealthHandler struct{}

type healthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(
	w http.ResponseWriter,
	r *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		healthResponse{
			Status:  "ok",
			Service: "auth",
		},
	)
}
