package logger

import (
	"encoding/json"
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
