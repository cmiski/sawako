package router

import (
	"net/http"

	"github.com/cmiski/sawako/gateway/internal/handlers"
	"github.com/cmiski/sawako/gateway/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	healthHandler *handlers.HealthHandler,
	authProxy http.Handler,
	eventProxy http.Handler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)

	r.Get("/health", healthHandler.Health)

	r.Handle("/auth/*", authProxy)
	r.Handle("/events/*", eventProxy)

	return r
}
