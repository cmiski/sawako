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
	"github.com/cmiski/sawako/gateway/internal/proxy"
	"github.com/cmiski/sawako/gateway/internal/router"
	"github.com/cmiski/sawako/gateway/internal/server"
)

func main() {
	cfg := config.Load()

	authProxy, err := proxy.NewReverseProxy(
		"auth",
		cfg.AuthServiceURL,
	)
	if err != nil {
		log.Fatalf(
			"failed to create auth proxy: %v",
			err,
		)
	}

	eventProxy, err := proxy.NewReverseProxy(
		"events",
		cfg.EventServiceURL,
	)
	if err != nil {
		log.Fatalf(
			"failed to create events proxy: %v",
			err,
		)
	}

	healthHandler := handlers.NewHealthHandler()

	r := router.NewRouter(
		healthHandler,
		authProxy,
		eventProxy,
	)

	srv := server.NewServer(cfg.Port, r)

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	log.Printf(
		"gateway listening on :%s (auth=%s, events=%s)",
		cfg.Port,
		cfg.AuthServiceURL,
		cfg.EventServiceURL,
	)

	sigCh := make(chan os.Signal, 1)

	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

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
