package main

import (
	"context"

	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

func main() {
	// Create base logger
	config := logger.DefaultConfig()
	config.Level = logger.INFO
	config.Formatter = formatter.NewJSONFormatter()
	config.Sinks = []logger.Sink{sink.NewConsoleSink()}
	baseLog := logger.NewWithConfig(config)
	defer baseLog.Close()

	// Create request-scoped logger with context
	requestLog := baseLog.With("request_id", "req_12345", "user_id", 42)

	// Store logger in context
	ctx := logger.WithContext(context.Background(), requestLog)

	// Simulate request processing
	handleRequest(ctx)
}

func handleRequest(ctx context.Context) {
	// Retrieve logger from context
	log := logger.FromContext(ctx)

	log.Info("processing_request")

	// Pass context to other functions
	authenticateUser(ctx)
	fetchUserData(ctx)

	log.Info("request_completed", "status", "success")
}

func authenticateUser(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Info("authenticating_user")
	// Authentication logic here
	log.Info("user_authenticated")
}

func fetchUserData(ctx context.Context) {
	log := logger.FromContext(ctx)
	log.Info("fetching_user_data", "table", "users")
	// Data fetching logic here
	log.Info("user_data_fetched", "records", 1)
}
