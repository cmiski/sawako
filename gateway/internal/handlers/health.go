package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// constructor for HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := HealthResponse{
		Status:  "ok",
		Service: "gateway",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
