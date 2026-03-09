package formatter

import (
	"encoding/json"
	"testing"

	"github.com/MohaCodez/structured-logger/logger"
)

func TestJSONFormatterBasic(t *testing.T) {
	formatter := NewJSONFormatter()
	entry := &logger.Entry{
		Timestamp: "2026-03-10T03:00:00+05:30",
		Level:     "INFO",
		Message:   "test_message",
		Fields:    map[string]interface{}{},
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if result["level"] != "INFO" {
		t.Errorf("expected level INFO, got %v", result["level"])
	}

	if result["message"] != "test_message" {
		t.Errorf("expected message test_message, got %v", result["message"])
	}
}

func TestJSONFormatterWithFields(t *testing.T) {
	formatter := NewJSONFormatter()
	entry := &logger.Entry{
		Timestamp: "2026-03-10T03:00:00+05:30",
		Level:     "ERROR",
		Message:   "error_occurred",
		Fields: map[string]interface{}{
			"user_id": 123,
			"error":   "connection timeout",
		},
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if result["user_id"] != float64(123) {
		t.Errorf("expected user_id 123, got %v", result["user_id"])
	}

	if result["error"] != "connection timeout" {
		t.Errorf("expected error message, got %v", result["error"])
	}
}

func TestJSONFormatterWithCaller(t *testing.T) {
	formatter := NewJSONFormatter()
	entry := &logger.Entry{
		Timestamp: "2026-03-10T03:00:00+05:30",
		Level:     "DEBUG",
		Message:   "debug_message",
		Caller:    "main.go:42",
		Fields:    map[string]interface{}{},
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("Format() error: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if result["caller"] != "main.go:42" {
		t.Errorf("expected caller main.go:42, got %v", result["caller"])
	}
}
