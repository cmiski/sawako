package router

import (
	"github.com/cmiski/sawako/gateway/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	healthHandler *handlers.HealthHandler, // dependency injection of the health handler
) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", healthHandler.Health)
	return r
}
