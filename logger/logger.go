package logger

import (
	"fmt"
	"os"

	"github.com/MohaCodez/structured-logger/async"
)

// Formatter is an interface that converts a log Entry into bytes for output.
// Implement this interface to create custom log formats (e.g., plain text, XML).
// The library provides a JSON formatter by default.
type Formatter interface {
	// Format takes a log entry and returns the formatted bytes to be written.
	// Returns an error if formatting fails.
	Format(entry *Entry) ([]byte, error)
}

// Sink is an interface for log output destinations.
// Implement this interface to send logs to custom destinations
// (e.g., databases, message queues, cloud services).
// The library provides console and file sinks by default.
type Sink interface {
	// Write outputs the formatted log data to the destination.
	// Returns an error if the write operation fails.
	Write(data []byte) error

	// Close releases any resources held by the sink (e.g., file handles).
	// Always call logger.Close() to ensure all sinks are properly closed.
	Close() error
}

// Logger is the main logging interface. It handles log entry creation,
// formatting, and dispatching to output sinks.
//
// Logger is safe for concurrent use by multiple goroutines.
// All fields are read-only after creation except when using With() which
// creates a new logger instance.
type Logger struct {
	level            Level                  // Minimum level to log (filters out lower levels)
	formatter        Formatter              // Converts entries to bytes
	sinks            []Sink                 // Output destinations (console, file, etc.)
	enableCaller     bool                   // Whether to capture file:line information
	asyncWorker      *async.Worker          // Background worker for async mode (nil if sync)
	contextFields    map[string]interface{} // Fields inherited by child loggers
	exitFunc         func(int)              // Function to call on Fatal (usually os.Exit)
	sinkErrorHandler func(error)            // Handler for sink write errors
}

// New creates a logger with the specified minimum log level and default settings.
// This is the simplest way to create a logger.
//
// Default settings:
//   - JSON formatter
//   - Console output
//   - Synchronous mode
//   - Caller tracing disabled
//
// Example:
//
//	log := logger.New(logger.INFO)
//	defer log.Close()
//	log.Info("application started")
//
// For more control, use NewWithConfig() instead.
func New(level Level) *Logger {
	config := DefaultConfig()
	config.Level = level
	return NewWithConfig(config)
}

// Close shuts down the logger gracefully.
// If async mode is enabled, this flushes the queue and waits for all pending
// logs to be written. Then it closes all sinks to release resources.
//
// Always defer Close() after creating a logger to ensure logs aren't lost:
//
//	log := logger.New(logger.INFO)
//	defer log.Close()
//
// Returns an error if any sink fails to close.
func (l *Logger) Close() error {
	// If async mode is enabled, stop the worker and flush the queue
	if l.asyncWorker != nil {
		l.asyncWorker.Stop() // Blocks until all queued logs are written
	}

	// Close all sinks to release resources (file handles, connections, etc.)
	for _, sink := range l.sinks {
		if err := sink.Close(); err != nil {
			return err
		}
	}
	return nil
}

// With creates a child logger that inherits all settings from the parent
// and adds additional context fields. The parent logger is not modified.
//
// Context fields are automatically included in every log from the child logger.
// This is useful for adding request IDs, user IDs, or other contextual information
// without passing them to every log call.
//
// Parameters are key-value pairs: With("key1", value1, "key2", value2, ...)
//
// Example:
//
//	// Base logger
//	baseLog := logger.New(logger.INFO)
//
//	// Request-scoped logger with context
//	requestLog := baseLog.With("request_id", "abc123", "user_id", 42)
//
//	// All logs from requestLog include request_id and user_id
//	requestLog.Info("processing request")
//	// Output: {"timestamp":"...","level":"INFO","message":"processing request","request_id":"abc123","user_id":42}
//
// Child loggers can be nested:
//
//	serviceLog := baseLog.With("service", "auth")
//	requestLog := serviceLog.With("request_id", "abc123")
//	// requestLog includes both "service" and "request_id"
func (l *Logger) With(keyValues ...interface{}) *Logger {
	// Parse the key-value pairs into a map
	fields := parseFields(keyValues)

	// Create new context fields map by merging parent fields with new fields
	// New fields override parent fields if keys conflict
	newContextFields := make(map[string]interface{}, len(l.contextFields)+len(fields))

	// Copy parent context fields
	for k, v := range l.contextFields {
		newContextFields[k] = v
	}

	// Add/override with new fields
	for k, v := range fields {
		newContextFields[k] = v
	}

	// Return a new logger instance with merged context
	// All other settings (level, formatter, sinks) are shared with parent
	return &Logger{
		level:            l.level,
		formatter:        l.formatter,
		sinks:            l.sinks,
		enableCaller:     l.enableCaller,
		asyncWorker:      l.asyncWorker, // Shared worker (safe for concurrent use)
		contextFields:    newContextFields,
		exitFunc:         l.exitFunc,
		sinkErrorHandler: l.sinkErrorHandler,
	}
}

// defaultFormatter is a simple JSON formatter used when no custom formatter is specified.
// It outputs basic log information in JSON format.
type defaultFormatter struct{}

// Format converts a log entry to a simple JSON string.
// This is a minimal implementation - use formatter.NewJSONFormatter() for full features.
func (f *defaultFormatter) Format(entry *Entry) ([]byte, error) {
	return []byte(fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s"}`,
		entry.Timestamp, entry.Level, entry.Message)), nil
}

// defaultConsoleSink writes logs to standard output (console).
// This is used when no custom sinks are specified.
type defaultConsoleSink struct{}

// Write outputs the log data to stdout (console).
func (s *defaultConsoleSink) Write(data []byte) error {
	fmt.Println(string(data))
	return nil
}

// Close does nothing for console output (no resources to release).
func (s *defaultConsoleSink) Close() error {
	return nil
}

// log is the internal method that handles all logging operations.
// All public log methods (Debug, Info, Warn, Error, Fatal) call this.
//
// Process:
//  1. Check if the log level meets the minimum threshold (fast filter)
//  2. Parse key-value pairs into a fields map
//  3. Merge context fields with call-specific fields
//  4. Create an Entry with timestamp, level, message, and fields
//  5. Format the entry to bytes
//  6. Write to all sinks (async or sync depending on configuration)
func (l *Logger) log(level Level, message string, keyValues ...interface{}) {
	// Early return if this log level is below the threshold
	// This is very fast (just an integer comparison) and avoids unnecessary work
	if level < l.level {
		return
	}

	// Parse the variadic key-value arguments into a map
	fields := parseFields(keyValues)

	// Merge context fields (from With()) with call-specific fields
	// Call-specific fields override context fields if keys conflict
	mergedFields := make(map[string]interface{}, len(l.contextFields)+len(fields))

	// First, copy all context fields
	for k, v := range l.contextFields {
		mergedFields[k] = v
	}

	// Then add/override with call-specific fields
	for k, v := range fields {
		mergedFields[k] = v
	}

	// Create the log entry with all metadata
	entry := newEntry(level, message, mergedFields, l.enableCaller)

	// Format the entry to bytes (JSON, text, etc.)
	data, err := l.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to format log entry: %v\n", err)
		return
	}

	// Write to sinks (output destinations)
	if l.asyncWorker != nil {
		// Async mode: enqueue for background processing
		// The log call returns immediately without waiting for I/O

		// Create a copy of sinks slice for the worker
		// (async.Sink is the same interface as logger.Sink)
		sinksCopy := make([]async.Sink, len(l.sinks))
		for i, s := range l.sinks {
			sinksCopy[i] = s
		}

		// Enqueue the log entry - returns immediately
		l.asyncWorker.Enqueue(data, sinksCopy)
	} else {
		// Sync mode: write to all sinks immediately (blocking)
		for _, sink := range l.sinks {
			if err := sink.Write(data); err != nil {
				// Don't fail the log call if one sink fails
				// Just report the error via the error handler
				l.sinkErrorHandler(err)
			}
		}
	}
}

// parseFields converts variadic key-value arguments into a map.
// Expects alternating keys (strings) and values (any type).
//
// Example: parseFields("user_id", 123, "name", "alice")
// Returns: map[string]interface{}{"user_id": 123, "name": "alice"}
//
// Handles edge cases:
//   - Odd number of arguments: adds "MISSING_VALUE" for the last key
//   - Non-string keys: skips the key-value pair and logs a warning
func parseFields(keyValues []interface{}) map[string]interface{} {
	fields := make(map[string]interface{})

	// Check if we have an odd number of arguments (missing a value)
	if len(keyValues)%2 != 0 {
		fmt.Fprintf(os.Stderr, "structured-logger: odd number of fields passed to log call, last key has no value\n")
		// Add a placeholder value so we don't panic
		keyValues = append(keyValues, "MISSING_VALUE")
	}

	// Process pairs: keyValues[0] and keyValues[1], then keyValues[2] and keyValues[3], etc.
	for i := 0; i < len(keyValues); i += 2 {
		// Keys must be strings
		key, ok := keyValues[i].(string)
		if !ok {
			// Skip this pair if the key isn't a string
			fmt.Fprintf(os.Stderr, "warning: non-string key at position %d, skipping pair\n", i)
			continue
		}
		// Values can be any type
		fields[key] = keyValues[i+1]
	}

	return fields
}

// Debug logs a message at DEBUG level with optional structured fields.
// DEBUG is the most verbose level, typically used for detailed troubleshooting.
// These logs are usually disabled in production.
//
// Parameters are key-value pairs: Debug("message", "key1", value1, "key2", value2, ...)
//
// Example:
//
//	log.Debug("processing request",
//	    "method", "GET",
//	    "path", "/api/users",
//	    "query_params", params,
//	)
func (l *Logger) Debug(message string, keyValues ...interface{}) {
	l.log(DEBUG, message, keyValues...)
}

// Info logs a message at INFO level with optional structured fields.
// INFO is for general informational messages about normal application operations.
// This is the recommended minimum level for production.
//
// Example:
//
//	log.Info("user logged in",
//	    "user_id", 123,
//	    "ip_address", "192.168.1.1",
//	)
func (l *Logger) Info(message string, keyValues ...interface{}) {
	l.log(INFO, message, keyValues...)
}

// Warn logs a message at WARN level with optional structured fields.
// WARN indicates potentially harmful situations that should be reviewed.
// The application continues to function normally.
//
// Example:
//
//	log.Warn("rate limit approaching",
//	    "user_id", 123,
//	    "current_requests", 95,
//	    "limit", 100,
//	)
func (l *Logger) Warn(message string, keyValues ...interface{}) {
	l.log(WARN, message, keyValues...)
}

// Error logs a message at ERROR level with optional structured fields.
// ERROR indicates error conditions that should be investigated.
// The application may continue but functionality may be impaired.
//
// Example:
//
//	log.Error("database connection failed",
//	    "error", err.Error(),
//	    "host", "db.example.com",
//	    "retry_count", 3,
//	)
func (l *Logger) Error(message string, keyValues ...interface{}) {
	l.log(ERROR, message, keyValues...)
}

// Fatal logs a message at FATAL level, then closes the logger and exits the program.
// FATAL indicates critical errors that require the application to terminate.
//
// Process:
//  1. Logs the message with FATAL level
//  2. Calls Close() to flush async queue and close sinks
//  3. Calls exitFunc (usually os.Exit(1)) to terminate the program
//
// Use Fatal only for unrecoverable errors. For recoverable errors, use Error instead.
//
// Example:
//
//	if err := initDatabase(); err != nil {
//	    log.Fatal("failed to initialize database",
//	        "error", err.Error(),
//	    )
//	    // Program exits here, code below won't execute
//	}
func (l *Logger) Fatal(message string, keyValues ...interface{}) {
	l.log(FATAL, message, keyValues...)
	l.Close()       // Flush logs and close sinks
	l.exitFunc(1)   // Exit the program with status code 1
}
