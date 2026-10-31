package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MeisterNickey/experience-hub/internal/config"
	"github.com/MeisterNickey/experience-hub/internal/repository"
	"github.com/MeisterNickey/experience-hub/internal/worker"
)

func main() {
	cfg := config.Load()

	postgresRepo := repository.NewPostgresRepository(cfg.PostgresDSN)
	dynamoRepo := repository.NewDynamoDBRepository(cfg.AWSRegion, cfg.DynamoDBTableName)
	redisRepo := repository.NewRedisRepository(cfg.RedisAddr, cfg.RedisPassword)

	processor := worker.NewProcessor(
		cfg.SQSQueueURL,
		cfg.WorkerConcurrency,
		postgresRepo,
		dynamoRepo,
		redisRepo,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("Experience-Hub Ticket Processor starting", "env", cfg.Environment)
	processor.Start(ctx)
}
