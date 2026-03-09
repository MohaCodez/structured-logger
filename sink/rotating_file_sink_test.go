package sink

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRotatingFileSinkBasic(t *testing.T) {
	testFile := "test_rotating.log"
	defer cleanupRotatingFiles(testFile)

	sink, err := NewRotatingFileSink(testFile, 1, 3)
	if err != nil {
		t.Fatalf("NewRotatingFileSink() error: %v", err)
	}
	defer sink.Close()

	// Write small log
	err = sink.Write([]byte("test log"))
	if err != nil {
		t.Errorf("Write() error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("log file was not created")
	}
}

func TestRotatingFileSinkRotation(t *testing.T) {
	testFile := "test_rotation.log"
	defer cleanupRotatingFiles(testFile)

	// Create sink with 1KB max size
	sink, err := NewRotatingFileSink(testFile, 0, 3) // 0 MB = very small for testing
	if err != nil {
		t.Fatalf("NewRotatingFileSink() error: %v", err)
	}
	defer sink.Close()

	// Write data to trigger rotation
	largeData := make([]byte, 2048)
	for i := range largeData {
		largeData[i] = 'A'
	}

	// First write
	sink.Write(largeData)

	// Second write should trigger rotation
	sink.Write(largeData)

	// Check that backup file exists
	backupFile := fmt.Sprintf("%s.1", testFile)
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Error("backup file was not created after rotation")
	}

	// Original file should still exist
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("original file should exist after rotation")
	}
}

func TestRotatingFileSinkMaxBackups(t *testing.T) {
	testFile := "test_max_backups.log"
	defer cleanupRotatingFiles(testFile)

	// Create sink with max 2 backups
	sink, err := NewRotatingFileSink(testFile, 0, 2)
	if err != nil {
		t.Fatalf("NewRotatingFileSink() error: %v", err)
	}
	defer sink.Close()

	largeData := make([]byte, 2048)
	for i := range largeData {
		largeData[i] = 'B'
	}

	// Trigger multiple rotations
	for i := 0; i < 5; i++ {
		sink.Write(largeData)
	}

	// Should have at most 2 backups
	backup1 := fmt.Sprintf("%s.1", testFile)
	backup2 := fmt.Sprintf("%s.2", testFile)
	backup3 := fmt.Sprintf("%s.3", testFile)

	if _, err := os.Stat(backup1); os.IsNotExist(err) {
		t.Error("backup .1 should exist")
	}

	if _, err := os.Stat(backup2); os.IsNotExist(err) {
		t.Error("backup .2 should exist")
	}

	// backup3 should not exist (exceeds maxBackups)
	if _, err := os.Stat(backup3); !os.IsNotExist(err) {
		t.Error("backup .3 should not exist (exceeds maxBackups)")
	}
}

func TestRotatingFileSinkConcurrent(t *testing.T) {
	testFile := "test_concurrent.log"
	defer cleanupRotatingFiles(testFile)

	sink, err := NewRotatingFileSink(testFile, 1, 3)
	if err != nil {
		t.Fatalf("NewRotatingFileSink() error: %v", err)
	}
	defer sink.Close()

	// Concurrent writes
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				sink.Write([]byte(fmt.Sprintf("log from goroutine %d iteration %d", id, j)))
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not crash or corrupt
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("log file should exist after concurrent writes")
	}
}

func cleanupRotatingFiles(basePath string) {
	os.Remove(basePath)
	for i := 1; i <= 10; i++ {
		os.Remove(fmt.Sprintf("%s.%d", basePath, i))
	}
	// Also clean up any test files in current directory
	matches, _ := filepath.Glob("test_*.log*")
	for _, match := range matches {
		os.Remove(match)
	}
}
