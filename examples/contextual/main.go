package main

import (
	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

func main() {
	// Create base logger
	config := logger.Config{
		Level:        logger.INFO,
		Formatter:    formatter.NewJSONFormatter(),
		Sinks:        []logger.Sink{sink.NewConsoleSink()},
		EnableCaller: true,
	}
	baseLog := logger.NewWithConfig(config)
	defer baseLog.Close()

	baseLog.Info("application_started", "version", "1.0.0")

	// Create service-level logger with context
	serviceLog := baseLog.With("service", "auth", "environment", "production")

	serviceLog.Info("service_initialized")

	// Simulate handling requests with request-specific context
	handleRequest(serviceLog, "req_001", 12345)
	handleRequest(serviceLog, "req_002", 67890)

	baseLog.Info("application_shutdown")
}

func handleRequest(serviceLog *logger.Logger, requestID string, userID int) {
	// Create request-specific logger
	reqLog := serviceLog.With("request_id", requestID, "user_id", userID)

	reqLog.Info("request_started")

	// Simulate processing
	authenticateUser(reqLog)
	fetchData(reqLog)

	reqLog.Info("request_completed", "status", "success")
}

func authenticateUser(reqLog *logger.Logger) {
	reqLog.Info("authenticating_user")
	// Authentication logic
	reqLog.Info("user_authenticated")
}

func fetchData(reqLog *logger.Logger) {
	reqLog.Info("fetching_data", "table", "users")
	// Data fetching logic
	reqLog.Info("data_fetched", "rows", 10)
}
