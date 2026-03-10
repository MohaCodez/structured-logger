package logger

import (
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// Entry represents a single log event with all its metadata.
// This struct holds all information about a log message before it's formatted and written.
type Entry struct {
	Timestamp string                 // When the log was created (RFC3339 format)
	Level     string                 // Severity level (DEBUG, INFO, WARN, ERROR, FATAL)
	Message   string                 // The main log message
	Caller    string                 // Source file and line number (e.g., "main.go:42")
	Fields    map[string]interface{} // Additional structured key-value data
}

// newEntry creates a new log entry with the provided information.
// This is an internal function called by the logger when a log method is invoked.
//
// Parameters:
//   - level: The severity level of this log entry
//   - message: The main log message
//   - fields: Additional structured data as key-value pairs
//   - enableCaller: Whether to capture the source file and line number
//
// Returns a fully populated Entry ready for formatting.
func newEntry(level Level, message string, fields map[string]interface{}, enableCaller bool) *Entry {
	entry := &Entry{
		Timestamp: time.Now().Format(time.RFC3339), // RFC3339 is a standard timestamp format
		Level:     level.String(),                  // Convert Level enum to string
		Message:   message,
		Fields:    fields,
	}

	// Optionally capture where in the code this log was called from
	if enableCaller {
		entry.Caller = getCaller()
	}

	return entry
}

// getCaller uses Go's runtime package to determine the file and line number
// where the log function was called. This helps developers quickly locate
// where a log message originated in their code.
//
// The number 4 in runtime.Caller(4) skips these stack frames:
//   0: runtime.Caller itself
//   1: getCaller (this function)
//   2: newEntry
//   3: the internal log() method
//   4: the public method (Debug/Info/Warn/Error/Fatal) - this is what we want
//
// Returns a string like "main.go:42" or "unknown" if caller info isn't available.
func getCaller() string {
	// runtime.Caller returns: program counter, file path, line number, and success flag
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		return "unknown"
	}
	// filepath.Base extracts just the filename from the full path
	// strconv.Itoa converts the line number (int) to a string
	return filepath.Base(file) + ":" + strconv.Itoa(line)
}
