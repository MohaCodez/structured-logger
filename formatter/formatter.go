package formatter

import "github.com/MohaCodez/structured-logger/logger"

// Formatter is an interface for converting log entries to bytes.
// Implement this interface to create custom output formats.
//
// The library provides JSONFormatter by default, but you can create
// formatters for plain text, XML, or any custom format.
//
// Example custom formatter:
//
//	type TextFormatter struct{}
//
//	func (f *TextFormatter) Format(entry *logger.Entry) ([]byte, error) {
//	    text := fmt.Sprintf("[%s] %s: %s\n",
//	        entry.Level, entry.Timestamp, entry.Message)
//	    return []byte(text), nil
//	}
type Formatter interface {
	// Format converts a log entry to bytes for output.
	// Returns the formatted bytes and any error that occurred during formatting.
	Format(entry *logger.Entry) ([]byte, error)
}
