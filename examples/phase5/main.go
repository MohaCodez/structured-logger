package main

import (
	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

func main() {
	// Setup formatter and sinks
	jsonFormatter := formatter.NewJSONFormatter()
	consoleSink := sink.NewConsoleSink()

	// Create logger with caller tracing enabled
	log := logger.NewWithCaller(
		logger.DEBUG,
		jsonFormatter,
		[]logger.Sink{consoleSink},
		true, // enable caller
	)
	defer log.Close()

	log.Info("application_started",
		"version", "2.0.0",
	)

	processRequest()
	handleError()

	log.Info("application_shutdown")
}

func processRequest() {
	log := getLogger()
	log.Debug("processing_request",
		"endpoint", "/api/users",
		"method", "GET",
	)
}

func handleError() {
	log := getLogger()
	log.Error("database_query_failed",
		"query", "SELECT * FROM users",
		"error", "connection timeout",
	)
}

func getLogger() *logger.Logger {
	jsonFormatter := formatter.NewJSONFormatter()
	consoleSink := sink.NewConsoleSink()
	return logger.NewWithCaller(
		logger.DEBUG,
		jsonFormatter,
		[]logger.Sink{consoleSink},
		true,
	)
}
