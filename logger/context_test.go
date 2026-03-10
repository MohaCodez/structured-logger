package logger

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func TestWithContext(t *testing.T) {
	mock := &mockSink{}
	config := DefaultConfig()
	config.Level = INFO
	config.Formatter = &jsonFormatterMock{}
	config.Sinks = []Sink{mock}
	log := NewWithConfig(config)

	ctx := WithContext(context.Background(), log)
	retrievedLog := FromContext(ctx)

	retrievedLog.Info("test_message", "key", "value")

	if len(mock.logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(mock.logs))
	}

	if !strings.Contains(mock.logs[0], "test_message") {
		t.Errorf("expected test_message in log, got: %s", mock.logs[0])
	}
}

func TestFromContextEmpty(t *testing.T) {
	// Empty context should return a non-nil default logger
	log := FromContext(context.Background())
	
	if log == nil {
		t.Error("FromContext should never return nil")
	}

	// Should be able to use the default logger
	log.Info("test")
}

// Mock sink that always returns an error
type errorSink struct {
	errorCount int
}

func (e *errorSink) Write(data []byte) error {
	e.errorCount++
	return fmt.Errorf("mock sink error")
}

func (e *errorSink) Close() error {
	return nil
}

func TestSinkErrorHandling(t *testing.T) {
	errorSink := &errorSink{}
	errorHandlerCalled := false
	var capturedError error

	config := DefaultConfig()
	config.Level = INFO
	config.Sinks = []Sink{errorSink}
	config.SinkErrorHandler = func(err error) {
		errorHandlerCalled = true
		capturedError = err
	}
	log := NewWithConfig(config)

	log.Info("test")

	if !errorHandlerCalled {
		t.Error("expected error handler to be called")
	}

	if capturedError == nil {
		t.Error("expected error to be captured")
	}

	if errorSink.errorCount != 1 {
		t.Errorf("expected 1 error, got %d", errorSink.errorCount)
	}
}
