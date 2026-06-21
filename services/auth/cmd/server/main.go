package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cmiski/sawako/services/auth/internal/authentication"
	"github.com/cmiski/sawako/services/auth/internal/config"
	"github.com/cmiski/sawako/services/auth/internal/handlers"
	"github.com/cmiski/sawako/services/auth/internal/jwt"
	"github.com/cmiski/sawako/services/auth/internal/postgres"
	"github.com/cmiski/sawako/services/auth/internal/project"
	"github.com/cmiski/sawako/services/auth/internal/refreshtoken"
	"github.com/cmiski/sawako/services/auth/internal/router"
	"github.com/cmiski/sawako/services/auth/internal/security"
	"github.com/cmiski/sawako/services/auth/internal/server"
	"github.com/cmiski/sawako/services/auth/internal/user"
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

	userRepo := postgres.NewUserRepository(pool)
	credentialRepo := postgres.NewCredentialRepository(pool)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(pool)
	projectRepo := postgres.NewProjectRepository(pool)
	txManager := postgres.NewTransactionManager(pool)

	userService := user.NewService(userRepo)
	refreshTokenService := refreshtoken.NewService()
	projectService := project.NewService(projectRepo)

	authService := authentication.NewService(
		userService,
		credentialRepo,
		refreshTokenRepo,
		refreshTokenService,
		security.NewBcryptPasswordHasher(),
		jwt.NewIssuer(
			cfg.JWTSecret,
			cfg.AccessTokenTTL,
			cfg.RefreshTokenTTL,
		),
		security.NewRandomRefreshTokenGenerator(),
		security.NewSHA256TokenHasher(),
		txManager,
	)

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	projectHandler := handlers.NewProjectHandler(projectService)

	r := router.NewRouter(
		healthHandler,
		authHandler,
		projectHandler,
	)

	srv := server.NewServer(cfg.Port, r)

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	log.Printf("auth service listening on :%s", cfg.Port)

	sigCh := make(chan os.Signal, 1)

	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("server gracefully stopped")
}
