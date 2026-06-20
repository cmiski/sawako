package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/cmiski/sawako/services/auth/internal/config"
	"github.com/cmiski/sawako/services/auth/internal/postgres"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	pool, err := postgres.NewPool(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatalf(
			"failed to connect to database: %v",
			err,
		)
	}
	defer pool.Close()

	migrationsDir := filepath.Join(
		"services",
		"auth",
		"migrations",
	)

	if err := postgres.RunMigrations(
		ctx,
		pool,
		migrationsDir,
	); err != nil {
		log.Fatalf(
			"failed to run migrations: %v",
			err,
		)
	}

	log.Println("migrations applied successfully")
}
