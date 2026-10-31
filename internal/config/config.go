package config

import (
	"os"
	"strconv"
)

// Config holds runtime configuration for both the Core API and SQS Worker.
type Config struct {
	Port              string
	Environment       string
	AWSRegion         string
	SQSQueueURL       string
	DynamoDBTableName string
	PostgresDSN       string
	RedisAddr         string
	RedisPassword     string
	WorkerConcurrency int
}

// Load reads configuration from environment variables with sensible local defaults.
func Load() *Config {
	concurrency, err := strconv.Atoi(getEnv("WORKER_CONCURRENCY", "5"))
	if err != nil || concurrency < 1 {
		concurrency = 5
	}

	return &Config{
		Port:              getEnv("PORT", "8080"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		AWSRegion:         getEnv("AWS_REGION", "us-east-1"),
		SQSQueueURL:       getEnv("SQS_QUEUE_URL", "http://localhost:4566/000000000000/experience-hub-tickets"),
		DynamoDBTableName: getEnv("DYNAMODB_TABLE_NAME", "experience-hub-audit-logs"),
		PostgresDSN:       getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/experience_hub?sslmode=disable"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		WorkerConcurrency: concurrency,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
