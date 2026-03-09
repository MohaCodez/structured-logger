package main

import (
	"time"

	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

func main() {
	// Setup formatter and sinks
	jsonFormatter := formatter.NewJSONFormatter()
	consoleSink := sink.NewConsoleSink()

	// Create async logger with buffer size 100
	log := logger.NewAsync(
		logger.INFO,
		jsonFormatter,
		[]logger.Sink{consoleSink},
		true,  // enable caller
		100,   // buffer size
	)
	defer log.Close()

	log.Info("async_logger_started")

	// Simulate high-throughput logging
	for i := 0; i < 10; i++ {
		log.Info("processing_request",
			"request_id", i,
			"user_id", 1000+i,
		)
		
		if i%3 == 0 {
			log.Warn("rate_limit_warning",
				"request_id", i,
				"threshold", 100,
			)
		}
	}

	log.Error("critical_error",
		"error", "database connection lost",
		"retry_count", 3,
	)

	log.Info("async_logger_shutdown")

	// Give async worker time to process queue before Close()
	time.Sleep(100 * time.Millisecond)
}
