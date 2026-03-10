package logger

import (
	"fmt"
	"os"

	"github.com/MohaCodez/structured-logger/async"
)

// BufferFullPolicy defines what happens when the async logging buffer is full.
// This only applies when async mode is enabled.
type BufferFullPolicy int

const (
	// BlockOnFull makes the logging call wait until buffer space is available.
	// This provides backpressure and ensures no logs are lost, but may slow down
	// the calling code if logs are being generated faster than they can be written.
	// This is the default and recommended for most applications.
	BlockOnFull BufferFullPolicy = iota

	// DropOnFull immediately drops the log entry if the buffer is full.
	// The calling code never blocks, but logs may be lost under high load.
	// Use this only if logging must never impact application performance.
	DropOnFull
)

// Config holds all configuration options for creating a logger.
// Use DefaultConfig() to get sensible defaults, then customize as needed.
type Config struct {
	// Level is the minimum severity level that will be logged.
	// Logs below this level are filtered out early (very fast).
	// Example: If set to INFO, DEBUG logs will be ignored.
	Level Level

	// Formatter converts log entries to bytes for output.
	// Default is JSON format. Implement the Formatter interface for custom formats.
	Formatter Formatter

	// Sinks are the output destinations (console, file, custom, etc.).
	// Logs are written to all sinks simultaneously (fan-out pattern).
	// Default is console output only.
	Sinks []Sink

	// EnableCaller adds file:line information to each log entry.
	// Useful for debugging but has a small performance cost.
	// Recommended: true in development, false in production.
	EnableCaller bool

	// Async enables non-blocking logging via a background goroutine.
	// Log calls return immediately after queuing the message.
	// Provides better throughput but logs may be lost if the process crashes
	// before the buffer is flushed. Always call Close() to flush on shutdown.
	Async bool

	// BufferSize is the queue size for async mode (ignored if Async is false).
	// Larger buffers handle bursts better but use more memory.
	// Typical values: 100-1000. Default: 100.
	BufferSize int

	// ExitFunc is called when Fatal() is invoked.
	// Defaults to os.Exit(1). Can be overridden for testing.
	ExitFunc func(int)

	// BufferFullPolicy determines behavior when async buffer is full.
	// Only applies when Async is true. Default: BlockOnFull.
	BufferFullPolicy BufferFullPolicy

	// SinkErrorHandler is called when a sink fails to write.
	// Defaults to printing to stderr. Customize to handle errors differently.
	SinkErrorHandler func(error)
}

// DefaultConfig returns a Config with sensible defaults for most applications.
// Customize the returned Config before passing it to NewWithConfig().
//
// Defaults:
//   - Level: INFO (production-appropriate)
//   - Formatter: JSON
//   - Sinks: Console only
//   - EnableCaller: false (better performance)
//   - Async: false (simpler, synchronous logging)
//   - BufferSize: 100
//   - ExitFunc: os.Exit
//   - BufferFullPolicy: BlockOnFull (no log loss)
func DefaultConfig() Config {
	return Config{
		Level:            INFO,
		Formatter:        &defaultFormatter{},
		Sinks:            []Sink{&defaultConsoleSink{}},
		EnableCaller:     false,
		Async:            false,
		BufferSize:       100,
		ExitFunc:         os.Exit,
		BufferFullPolicy: BlockOnFull,
	}
}

// NewWithConfig creates a new logger with the specified configuration.
// This is the recommended way to create a logger with custom settings.
//
// Example:
//
//	config := logger.DefaultConfig()
//	config.Level = logger.DEBUG
//	config.EnableCaller = true
//	log := logger.NewWithConfig(config)
//	defer log.Close() // Important: flush async buffer and close sinks
func NewWithConfig(config Config) *Logger {
	// Use provided ExitFunc or default to os.Exit
	exitFunc := config.ExitFunc
	if exitFunc == nil {
		exitFunc = os.Exit
	}

	// Use provided error handler or default to stderr output
	sinkErrorHandler := config.SinkErrorHandler
	if sinkErrorHandler == nil {
		sinkErrorHandler = func(err error) {
			fmt.Fprintf(os.Stderr, "sink write error: %v\n", err)
		}
	}

	// Create the logger with the specified configuration
	logger := &Logger{
		level:            config.Level,
		formatter:        config.Formatter,
		sinks:            config.Sinks,
		enableCaller:     config.EnableCaller,
		asyncWorker:      nil, // Will be set below if async is enabled
		contextFields:    make(map[string]interface{}),
		exitFunc:         exitFunc,
		sinkErrorHandler: sinkErrorHandler,
	}

	// If async mode is enabled, create a background worker goroutine
	if config.Async {
		// The worker handles writing logs in the background
		// dropOnFull is true when BufferFullPolicy is DropOnFull
		logger.asyncWorker = async.NewWorker(
			config.BufferSize,
			config.BufferFullPolicy == DropOnFull,
			sinkErrorHandler,
		)
	}

	return logger
}
