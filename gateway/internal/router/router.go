package router

import (
	"github.com/cmiski/sawako/gateway/internal/handlers"
	"github.com/cmiski/sawako/gateway/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	healthHandler *handlers.HealthHandler, // dependency injection of the health handler
) *chi.Mux {
	r := chi.NewRouter()

	// Add logging middleware to log incoming requests and their latency
	r.Use(middleware.Logging)
	// Add recovery middleware to handle panics gracefully
	r.Use(middleware.Recovery)

	r.Get("/health", healthHandler.Health)
	return r
}
