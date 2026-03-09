package main

import (
	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

func main() {
	// Setup formatter
	jsonFormatter := formatter.NewJSONFormatter()

	// Setup sinks
	consoleSink := sink.NewConsoleSink()
	fileSink, err := sink.NewFileSink("app.log")
	if err != nil {
		panic(err)
	}
	defer fileSink.Close()

	// Create logger with multiple sinks
	sinks := []logger.Sink{consoleSink, fileSink}
	log := logger.NewWithSinks(logger.INFO, jsonFormatter, sinks)
	defer log.Close()

	log.Info("application_started",
		"version", "1.0.0",
		"environment", "production",
	)

	log.Warn("disk_space_low",
		"available_gb", 5.2,
		"threshold_gb", 10.0,
	)

	log.Error("payment_processing_failed",
		"transaction_id", "txn_12345",
		"amount", 99.99,
		"error", "gateway timeout",
	)

	log.Info("application_shutdown")
}
