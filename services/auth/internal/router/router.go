package router

import (
	"github.com/cmiski/sawako/services/auth/internal/handlers"
	"github.com/cmiski/sawako/services/auth/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recovery)

	r.Get("/health", healthHandler.Health)
	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/refresh", authHandler.Refresh)

	return r
}
