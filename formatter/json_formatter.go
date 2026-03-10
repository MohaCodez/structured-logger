package formatter

import (
	"encoding/json"

	"github.com/MohaCodez/structured-logger/logger"
)

// JSONFormatter formats log entries as JSON objects.
// This is the recommended formatter for production use because:
//   - JSON is machine-readable (easy to parse by log aggregation tools)
//   - Structured fields are preserved with their types
//   - Widely supported by logging infrastructure
//
// Output format:
//
//	{
//	  "timestamp": "2026-03-11T03:00:00+05:30",
//	  "level": "INFO",
//	  "message": "user logged in",
//	  "caller": "main.go:42",        // if caller tracing is enabled
//	  "user_id": 123,                // structured fields
//	  "ip_address": "192.168.1.1"
//	}
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter.
// This is the standard way to create a JSONFormatter instance.
//
// Example:
//
//	config := logger.DefaultConfig()
//	config.Formatter = formatter.NewJSONFormatter()
//	log := logger.NewWithConfig(config)
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format converts a log entry to JSON bytes.
//
// Process:
//  1. Create a map with standard fields (timestamp, level, message)
//  2. Add caller field if present
//  3. Add all structured fields from the entry
//  4. Marshal the map to JSON bytes
//
// Returns the JSON bytes or an error if marshaling fails.
func (f *JSONFormatter) Format(entry *logger.Entry) ([]byte, error) {
	// Create a map to hold all fields
	// Using interface{} as value type allows any data type
	m := map[string]interface{}{
		"timestamp": entry.Timestamp,
		"level":     entry.Level,
		"message":   entry.Message,
	}

	// Add caller information if it was captured
	if entry.Caller != "" {
		m["caller"] = entry.Caller
	}

	// Add all structured fields to the map
	// These fields come from log calls like: log.Info("msg", "key", value)
	for k, v := range entry.Fields {
		m[k] = v
	}

	// Convert the map to JSON bytes
	// json.Marshal handles all Go types (strings, numbers, bools, slices, maps, etc.)
	return json.Marshal(m)
}
