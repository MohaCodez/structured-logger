package main

import (
	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
)

func main() {
	// Create logger with JSON formatter
	jsonFormatter := formatter.NewJSONFormatter()
	log := logger.NewWithFormatter(logger.INFO, jsonFormatter)

	log.Info("service_started",
		"port", 8080,
		"environment", "production",
	)

	log.Warn("high_cpu_usage",
		"cpu_percent", 85.5,
		"threshold", 80.0,
	)

	log.Error("request_failed",
		"endpoint", "/api/users",
		"status_code", 500,
		"error", "internal server error",
	)
}
