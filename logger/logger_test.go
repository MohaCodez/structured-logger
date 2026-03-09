package logger

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// Mock sink for testing
type mockSink struct {
	logs []string
}

func (m *mockSink) Write(data []byte) error {
	m.logs = append(m.logs, string(data))
	return nil
}

func (m *mockSink) Close() error {
	return nil
}

func TestLoggerLevelFiltering(t *testing.T) {
	mock := &mockSink{}
	log := NewWithSinks(WARN, &defaultFormatter{}, []Sink{mock})

	log.Debug("debug")
	log.Info("info")
	log.Warn("warn")
	log.Error("error")

	if len(mock.logs) != 2 {
		t.Errorf("expected 2 logs, got %d", len(mock.logs))
	}

	if !strings.Contains(mock.logs[0], "WARN") {
		t.Errorf("expected WARN log, got %s", mock.logs[0])
	}

	if !strings.Contains(mock.logs[1], "ERROR") {
		t.Errorf("expected ERROR log, got %s", mock.logs[1])
	}
}

func TestLoggerWithFields(t *testing.T) {
	mock := &mockSink{}
	log := NewWithSinks(INFO, &defaultFormatter{}, []Sink{mock})

	log.Info("test", "key1", "value1", "key2", 123)

	if len(mock.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(mock.logs))
	}

	// Parse JSON to verify fields
	var result map[string]interface{}
	err := json.Unmarshal([]byte(mock.logs[0]), &result)
	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if result["message"] != "test" {
		t.Errorf("expected message 'test', got %v", result["message"])
	}
}

func TestLoggerMultipleSinks(t *testing.T) {
	mock1 := &mockSink{}
	mock2 := &mockSink{}
	log := NewWithSinks(INFO, &defaultFormatter{}, []Sink{mock1, mock2})

	log.Info("test")

	if len(mock1.logs) != 1 {
		t.Errorf("sink1: expected 1 log, got %d", len(mock1.logs))
	}

	if len(mock2.logs) != 1 {
		t.Errorf("sink2: expected 1 log, got %d", len(mock2.logs))
	}
}

func TestLoggerWithCaller(t *testing.T) {
	mock := &mockSink{}
	// Use JSONFormatter which supports caller field
	jsonFormatter := &jsonFormatterMock{}
	log := NewWithCaller(INFO, jsonFormatter, []Sink{mock}, true)

	log.Info("test")

	if len(mock.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(mock.logs))
	}

	if !strings.Contains(mock.logs[0], "logger_test.go") {
		t.Errorf("expected caller field with logger_test.go in log, got: %s", mock.logs[0])
	}
}

// Mock JSON formatter for testing
type jsonFormatterMock struct{}

func (f *jsonFormatterMock) Format(entry *Entry) ([]byte, error) {
	result := `{"level":"` + entry.Level + `","message":"` + entry.Message + `"`
	if entry.Caller != "" {
		result += `,"caller":"` + entry.Caller + `"`
	}
	// Add fields
	for k, v := range entry.Fields {
		result += fmt.Sprintf(`,"%s":`, k)
		switch val := v.(type) {
		case string:
			result += fmt.Sprintf(`"%s"`, val)
		case int:
			result += fmt.Sprintf(`%d`, val)
		default:
			result += fmt.Sprintf(`"%v"`, val)
		}
	}
	result += `}`
	return []byte(result), nil
}

func TestLoggerClose(t *testing.T) {
	mock := &mockSink{}
	log := NewWithSinks(INFO, &defaultFormatter{}, []Sink{mock})

	err := log.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}

func TestLoggerWith(t *testing.T) {
	mock := &mockSink{}
	baseLog := NewWithSinks(INFO, &jsonFormatterMock{}, []Sink{mock})

	// Create child logger with context fields
	childLog := baseLog.With("request_id", "abc123", "user_id", 42)

	childLog.Info("test_message")

	if len(mock.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(mock.logs))
	}

	log := mock.logs[0]
	if !strings.Contains(log, "request_id") || !strings.Contains(log, "abc123") {
		t.Errorf("expected request_id in log, got: %s", log)
	}
	if !strings.Contains(log, "user_id") {
		t.Errorf("expected user_id in log, got: %s", log)
	}
}

func TestLoggerWithParentUnchanged(t *testing.T) {
	mock := &mockSink{}
	baseLog := NewWithSinks(INFO, &jsonFormatterMock{}, []Sink{mock})

	// Create child logger
	childLog := baseLog.With("request_id", "abc123")

	// Log from parent - should not have request_id
	baseLog.Info("parent_message")

	if len(mock.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(mock.logs))
	}

	if strings.Contains(mock.logs[0], "request_id") {
		t.Errorf("parent logger should not have child context fields")
	}

	// Log from child - should have request_id
	childLog.Info("child_message")

	if len(mock.logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(mock.logs))
	}

	if !strings.Contains(mock.logs[1], "request_id") {
		t.Errorf("child logger should have context fields")
	}
}

func TestLoggerWithNested(t *testing.T) {
	mock := &mockSink{}
	baseLog := NewWithSinks(INFO, &jsonFormatterMock{}, []Sink{mock})

	// Create nested child loggers
	serviceLog := baseLog.With("service", "auth")
	requestLog := serviceLog.With("request_id", "123")

	requestLog.Info("test")

	if len(mock.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(mock.logs))
	}

	log := mock.logs[0]
	if !strings.Contains(log, "service") || !strings.Contains(log, "auth") {
		t.Errorf("expected service field from parent, got: %s", log)
	}
	if !strings.Contains(log, "request_id") || !strings.Contains(log, "123") {
		t.Errorf("expected request_id field, got: %s", log)
	}
}

func TestLoggerWithFieldOverride(t *testing.T) {
	mock := &mockSink{}
	baseLog := NewWithSinks(INFO, &jsonFormatterMock{}, []Sink{mock})

	// Create child with context field
	childLog := baseLog.With("key", "context_value")

	// Log with same key - should override
	childLog.Info("test", "key", "call_value")

	if len(mock.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(mock.logs))
	}

	// Call value should override context value
	if !strings.Contains(mock.logs[0], "call_value") {
		t.Errorf("call fields should override context fields, got: %s", mock.logs[0])
	}
}
