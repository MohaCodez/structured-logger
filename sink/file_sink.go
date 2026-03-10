package sink

import (
	"fmt"
	"os"
)

// FileSink writes log output to a file on disk.
// Logs are appended to the file, so existing content is preserved.
//
// Use FileSink when you need:
//   - Persistent logs that survive application restarts
//   - Log files for auditing or compliance
//   - Logs that can be analyzed later
//
// For automatic log rotation, use RotatingFileSink instead.
//
// Note: FileSink keeps the file open for the lifetime of the logger.
// Always call logger.Close() to ensure the file is properly closed.
type FileSink struct {
	file *os.File // The open file handle for writing logs
}

// NewFileSink creates a new file sink that writes to the specified path.
// The file is opened in append mode, so existing content is preserved.
//
// File permissions: 0644 (owner can read/write, others can read)
//
// Returns an error if:
//   - The file cannot be created (e.g., invalid path, no permissions)
//   - The directory doesn't exist
//   - The disk is full
//
// Example:
//
//	fileSink, err := sink.NewFileSink("/var/log/myapp/app.log")
//	if err != nil {
//	    log.Fatal("failed to create file sink", "error", err)
//	}
//
//	config := logger.DefaultConfig()
//	config.Sinks = []logger.Sink{fileSink}
//	log := logger.NewWithConfig(config)
//	defer log.Close() // Important: closes the file
func NewFileSink(path string) (*FileSink, error) {
	// Open the file with these flags:
	//   os.O_CREATE: Create the file if it doesn't exist
	//   os.O_WRONLY: Open for writing only (not reading)
	//   os.O_APPEND: Write to the end of the file (don't overwrite)
	// Permission 0644: rw-r--r-- (owner read/write, group/others read)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		// %w wraps the error so callers can use errors.Is() and errors.As()
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return &FileSink{file: file}, nil
}

// Write appends the log data to the file.
// Each log entry is written as a single line.
//
// Returns an error if:
//   - The disk is full
//   - The file was closed or deleted
//   - File permissions changed
func (s *FileSink) Write(data []byte) error {
	// fmt.Fprintln writes to the file and adds a newline
	// The underscore (_) discards the number of bytes written
	_, err := fmt.Fprintln(s.file, string(data))
	return err
}

// Close closes the file handle and releases the file descriptor.
// After calling Close, no more writes can be performed.
//
// This is automatically called by logger.Close().
// Always defer logger.Close() to ensure files are properly closed.
func (s *FileSink) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
