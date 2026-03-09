package main

import (
	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
)

func main() {
	// Start with default and customize
	config := logger.DefaultConfig()
	config.Level = logger.DEBUG
	config.EnableCaller = true
	config.Formatter = formatter.NewJSONFormatter()

	log := logger.NewWithConfig(config)
	defer log.Close()

	log.Debug("application_started",
		"version", "3.0.0",
		"environment", "development",
	)

	log.Info("server_listening",
		"port", 8080,
		"protocol", "http",
	)

	log.Warn("deprecated_api_used",
		"endpoint", "/api/v1/users",
		"replacement", "/api/v2/users",
	)

	log.Error("authentication_failed",
		"user", "alice",
		"reason", "invalid_token",
	)
}
