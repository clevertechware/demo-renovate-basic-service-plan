package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clevertechware/demo-renovate-basic-service-plan/internal/api"
	"github.com/clevertechware/demo-renovate-basic-service-plan/internal/plan"
)

const (
	defaultDatabaseURL = "postgres://plan:plan@localhost:5432/plan?sslmode=disable"
	listenAddr         = ":8081"
	shutdownTimeout    = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("service-plan stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	repo := plan.Open(envOr("DATABASE_URL", defaultDatabaseURL))
	defer func() {
		if err := repo.Close(); err != nil {
			slog.Error("close database", "error", err)
		}
	}()

	if err := repo.Migrate(ctx); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	if err := repo.Seed(ctx, plan.DemoPlans); err != nil {
		return fmt.Errorf("seed: %w", err)
	}

	server := &http.Server{
		Addr:              listenAddr,
		Handler:           api.NewRouter(repo),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return serve(ctx, server)
}

func serve(ctx context.Context, server *http.Server) error {
	errc := make(chan error, 1)
	go func() { errc <- server.ListenAndServe() }()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
