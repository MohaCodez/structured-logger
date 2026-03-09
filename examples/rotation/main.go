package main

import (
	"fmt"
	"os"

	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

func main() {
	// Clean up any existing test files
	os.Remove("rotating.log")
	os.Remove("rotating.log.1")
	os.Remove("rotating.log.2")
	os.Remove("rotating.log.3")

	// Create rotating file sink
	// MaxSize: 1KB (very small for demo purposes)
	// MaxBackups: 3
	rotatingSink, err := sink.NewRotatingFileSink("rotating.log", 0, 3)
	if err != nil {
		panic(err)
	}
	defer rotatingSink.Close()

	// Create logger with rotating sink
	config := logger.Config{
		Level:     logger.INFO,
		Formatter: formatter.NewJSONFormatter(),
		Sinks:     []logger.Sink{sink.NewConsoleSink(), rotatingSink},
	}

	log := logger.NewWithConfig(config)
	defer log.Close()

	log.Info("rotation_demo_started")

	// Generate enough logs to trigger rotation
	for i := 0; i < 50; i++ {
		log.Info("log_entry",
			"iteration", i,
			"data", "This is a log message with some data to increase size",
			"more_data", "Additional padding to make the log entry larger",
		)
	}

	log.Info("rotation_demo_completed")

	// Check which files were created
	fmt.Println("\n=== Files created ===")
	files := []string{"rotating.log", "rotating.log.1", "rotating.log.2", "rotating.log.3"}
	for _, file := range files {
		if info, err := os.Stat(file); err == nil {
			fmt.Printf("%s: %d bytes\n", file, info.Size())
		}
	}
}
