package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MeisterNickey/experience-hub/internal/config"
	"github.com/MeisterNickey/experience-hub/internal/handlers"
	"github.com/MeisterNickey/experience-hub/internal/repository"
)

func main() {
	cfg := config.Load()

	postgresRepo := repository.NewPostgresRepository(cfg.PostgresDSN)
	redisRepo := repository.NewRedisRepository(cfg.RedisAddr, cfg.RedisPassword)
	ticketHandler := handlers.NewTicketHandler(postgresRepo, redisRepo, cfg.SQSQueueURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", ticketHandler.HealthCheck)
	mux.HandleFunc("/api/v1/tickets", ticketHandler.IngestTicket)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Experience-Hub Core API listening", "port", cfg.Port, "env", cfg.Environment)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("shutting down Core API server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
