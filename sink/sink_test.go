package sink

import (
	"os"
	"testing"
)

func TestConsoleSink(t *testing.T) {
	sink := NewConsoleSink()

	err := sink.Write([]byte("test log"))
	if err != nil {
		t.Errorf("Write() error: %v", err)
	}

	err = sink.Close()
	if err != nil {
		t.Errorf("Close() error: %v", err)
	}
}

func TestFileSink(t *testing.T) {
	testFile := "test_sink.log"
	defer os.Remove(testFile)

	sink, err := NewFileSink(testFile)
	if err != nil {
		t.Fatalf("NewFileSink() error: %v", err)
	}

	err = sink.Write([]byte("test log line 1"))
	if err != nil {
		t.Errorf("Write() error: %v", err)
	}

	err = sink.Write([]byte("test log line 2"))
	if err != nil {
		t.Errorf("Write() error: %v", err)
	}

	err = sink.Close()
	if err != nil {
		t.Errorf("Close() error: %v", err)
	}

	// Verify file contents
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("file is empty")
	}
}

func TestFileSinkInvalidPath(t *testing.T) {
	_, err := NewFileSink("/invalid/path/test.log")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
