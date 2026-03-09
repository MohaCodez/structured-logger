package logger

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Level != INFO {
		t.Errorf("expected default level INFO, got %v", config.Level)
	}

	if config.EnableCaller {
		t.Error("expected EnableCaller to be false by default")
	}

	if config.Async {
		t.Error("expected Async to be false by default")
	}

	if config.BufferSize != 100 {
		t.Errorf("expected default BufferSize 100, got %d", config.BufferSize)
	}

	if config.Formatter == nil {
		t.Error("expected default formatter to be set")
	}

	if len(config.Sinks) == 0 {
		t.Error("expected default sinks to be set")
	}
}

func TestNewWithConfig(t *testing.T) {
	mock := &mockSink{}
	config := Config{
		Level:        DEBUG,
		Formatter:    &defaultFormatter{},
		Sinks:        []Sink{mock},
		EnableCaller: true,
		Async:        false,
		BufferSize:   50,
	}

	log := NewWithConfig(config)
	defer log.Close()

	if log.level != DEBUG {
		t.Errorf("expected level DEBUG, got %v", log.level)
	}

	if !log.enableCaller {
		t.Error("expected enableCaller to be true")
	}

	if log.asyncWorker != nil {
		t.Error("expected asyncWorker to be nil for sync mode")
	}
}

func TestNewWithConfigAsync(t *testing.T) {
	mock := &mockSink{}
	config := Config{
		Level:        INFO,
		Formatter:    &defaultFormatter{},
		Sinks:        []Sink{mock},
		EnableCaller: false,
		Async:        true,
		BufferSize:   200,
	}

	log := NewWithConfig(config)
	defer log.Close()

	if log.asyncWorker == nil {
		t.Error("expected asyncWorker to be initialized for async mode")
	}
}

func TestConfigCustomization(t *testing.T) {
	config := DefaultConfig()
	config.Level = ERROR
	config.EnableCaller = true

	log := NewWithConfig(config)
	defer log.Close()

	if log.level != ERROR {
		t.Errorf("expected level ERROR, got %v", log.level)
	}

	if !log.enableCaller {
		t.Error("expected enableCaller to be true")
	}
}
