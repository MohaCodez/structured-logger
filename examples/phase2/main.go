package main

import "github.com/MohaCodez/structured-logger/logger"

func main() {
	log := logger.New(logger.DEBUG)

	// Simple message without fields
	log.Info("service started")

	// Message with structured fields
	log.Info("user_login",
		"user_id", 123,
		"username", "alice",
		"ip", "10.1.2.4",
	)

	log.Warn("high_memory_usage",
		"memory_mb", 1024,
		"threshold_mb", 800,
		"service", "api-server",
	)

	log.Error("database_connection_failed",
		"host", "db.example.com",
		"port", 5432,
		"retry_count", 3,
		"error", "connection timeout",
	)

	// Mixed types
	log.Debug("request_processed",
		"duration_ms", 45.3,
		"success", true,
		"endpoint", "/api/users",
	)
}
