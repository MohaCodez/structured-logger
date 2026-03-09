package main

import (
	"time"

	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

// Complete example demonstrating all logger features
func main() {
	// Setup: Create logger with full configuration
	fileSink, err := sink.NewFileSink("complete_example.log")
	if err != nil {
		panic(err)
	}

	config := logger.Config{
		Level:        logger.DEBUG,
		Formatter:    formatter.NewJSONFormatter(),
		Sinks:        []logger.Sink{sink.NewConsoleSink(), fileSink},
		EnableCaller: true,
		Async:        true,
		BufferSize:   200,
	}

	log := logger.NewWithConfig(config)
	defer log.Close()

	// 1. Basic logging at different levels
	log.Debug("application_initializing")
	log.Info("application_started", "version", "1.0.0", "environment", "production")

	// 2. Structured fields with various types
	log.Info("server_configuration",
		"host", "0.0.0.0",
		"port", 8080,
		"ssl_enabled", true,
		"max_connections", 1000,
		"timeout_seconds", 30.5,
	)

	// 3. Simulated request processing
	for i := 1; i <= 3; i++ {
		processRequest(log, i)
	}

	// 4. Warning conditions
	log.Warn("memory_usage_high",
		"current_mb", 850,
		"threshold_mb", 800,
		"percentage", 85.5,
	)

	// 5. Error handling
	log.Error("database_connection_failed",
		"host", "db.example.com",
		"port", 5432,
		"error", "connection timeout",
		"retry_count", 3,
	)

	// 6. Performance metrics
	start := time.Now()
	performOperation()
	log.Info("operation_completed",
		"operation", "data_sync",
		"duration_ms", time.Since(start).Milliseconds(),
		"records_processed", 1500,
	)

	// 7. User activity
	log.Info("user_action",
		"user_id", 12345,
		"username", "alice",
		"action", "file_upload",
		"file_size_mb", 2.5,
		"success", true,
	)

	// 8. System events
	log.Info("cache_cleared",
		"cache_type", "redis",
		"keys_removed", 250,
	)

	log.Info("application_shutdown", "uptime_seconds", 3600)
}

func processRequest(log *logger.Logger, requestID int) {
	log.Info("request_received",
		"request_id", requestID,
		"method", "POST",
		"path", "/api/users",
		"client_ip", "192.168.1.100",
	)

	// Simulate processing
	time.Sleep(10 * time.Millisecond)

	if requestID%2 == 0 {
		log.Warn("rate_limit_warning",
			"request_id", requestID,
			"user_id", 100+requestID,
			"requests_per_minute", 120,
			"limit", 100,
		)
	}

	log.Info("request_completed",
		"request_id", requestID,
		"status_code", 200,
		"response_time_ms", 45,
	)
}

func performOperation() {
	time.Sleep(50 * time.Millisecond)
}
