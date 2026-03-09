package main

import (
	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

func main() {
	// Example 1: Using default config
	config := logger.DefaultConfig()
	log := logger.NewWithConfig(config)
	defer log.Close()

	log.Info("using_default_config")

	// Example 2: Custom synchronous config
	syncConfig := logger.Config{
		Level:        logger.DEBUG,
		Formatter:    formatter.NewJSONFormatter(),
		Sinks:        []logger.Sink{sink.NewConsoleSink()},
		EnableCaller: true,
		Async:        false,
	}
	syncLog := logger.NewWithConfig(syncConfig)
	defer syncLog.Close()

	syncLog.Debug("sync_logger_with_caller",
		"mode", "synchronous",
		"caller_enabled", true,
	)

	// Example 3: Custom async config with file sink
	fileSink, err := sink.NewFileSink("config_example.log")
	if err != nil {
		panic(err)
	}

	asyncConfig := logger.Config{
		Level:        logger.INFO,
		Formatter:    formatter.NewJSONFormatter(),
		Sinks:        []logger.Sink{sink.NewConsoleSink(), fileSink},
		EnableCaller: true,
		Async:        true,
		BufferSize:   200,
	}
	asyncLog := logger.NewWithConfig(asyncConfig)
	defer asyncLog.Close()

	asyncLog.Info("async_logger_started",
		"mode", "asynchronous",
		"buffer_size", 200,
	)

	for i := 0; i < 5; i++ {
		asyncLog.Info("processing_item",
			"item_id", i,
			"status", "completed",
		)
	}

	asyncLog.Warn("high_load_detected",
		"queue_size", 150,
		"threshold", 100,
	)
}
