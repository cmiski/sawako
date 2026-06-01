package main

import (
	"log"

	"github.com/cmiski/sawako/gateway/internal/config"
	"github.com/cmiski/sawako/gateway/internal/handlers"
	"github.com/cmiski/sawako/gateway/internal/router"
	"github.com/cmiski/sawako/gateway/internal/server"
)

func main() {
	cfg := config.Load()

	healthHandler := handlers.NewHealthHandler()

	r := router.NewRouter(
		healthHandler,
	)

	srv := server.NewServer(cfg.Port, r)

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
