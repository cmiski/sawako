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

	// Start the server in a separate goroutine
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// channel to listen for interrupt or terminate signals
	sigCh := make(chan os.Signal, 1)

	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	// Block until we receive a signal
	<-sigCh

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("server gracefully stopped")
}
