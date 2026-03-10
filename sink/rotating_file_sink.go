package sink

import (
	"fmt"
	"os"
	"sync"
)

// RotatingFileSink writes logs to a file and automatically rotates it
// when it reaches a maximum size. This prevents log files from growing
// indefinitely and consuming all disk space.
//
// How rotation works:
//  1. When app.log reaches maxSizeMB, it's renamed to app.log.1
//  2. Previous backups shift: app.log.1 → app.log.2, app.log.2 → app.log.3, etc.
//  3. A new empty app.log is created
//  4. Oldest backups beyond maxBackups are deleted
//
// Example with maxBackups=3:
//   app.log       (current, active file)
//   app.log.1     (most recent backup)
//   app.log.2     (older backup)
//   app.log.3     (oldest backup)
//
// RotatingFileSink is thread-safe and can be used by multiple goroutines.
type RotatingFileSink struct {
	path       string      // Path to the log file (e.g., "/var/log/app.log")
	maxSizeMB  int64       // Maximum file size in bytes before rotation
	maxBackups int         // Maximum number of backup files to keep
	file       *os.File    // Current open file handle
	size       int64       // Current size of the file in bytes
	mu         sync.Mutex  // Mutex to protect concurrent writes and rotation
}

// NewRotatingFileSink creates a new rotating file sink.
//
// Parameters:
//   - path: Path to the log file (e.g., "/var/log/myapp/app.log")
//   - maxSizeMB: Maximum file size in megabytes before rotation (e.g., 10 for 10MB)
//   - maxBackups: Maximum number of backup files to keep (e.g., 5 keeps app.log.1 through app.log.5)
//
// Returns an error if the file cannot be opened or its size cannot be determined.
//
// Example:
//
//	// Rotate when file reaches 10MB, keep 5 backups
//	rotatingSink, err := sink.NewRotatingFileSink("/var/log/app.log", 10, 5)
//	if err != nil {
//	    log.Fatal("failed to create rotating sink", "error", err)
//	}
//
//	config := logger.DefaultConfig()
//	config.Sinks = []logger.Sink{rotatingSink}
//	log := logger.NewWithConfig(config)
//	defer log.Close()
func NewRotatingFileSink(path string, maxSizeMB int, maxBackups int) (*RotatingFileSink, error) {
	// Open the file in append mode (create if doesn't exist)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Get the current file size so we know when to rotate
	info, err := file.Stat()
	if err != nil {
		file.Close() // Clean up the file handle before returning error
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &RotatingFileSink{
		path:       path,
		maxSizeMB:  int64(maxSizeMB) * 1024 * 1024, // Convert MB to bytes
		maxBackups: maxBackups,
		file:       file,
		size:       info.Size(), // Track current size
	}, nil
}

// Write appends log data to the file and rotates if necessary.
// This method is thread-safe - multiple goroutines can call it concurrently.
//
// Process:
//  1. Lock the mutex to prevent concurrent writes/rotations
//  2. Check if adding this log would exceed maxSizeMB
//  3. If yes, rotate the file (rename current, create new)
//  4. Write the log data to the file
//  5. Update the size counter
//  6. Unlock the mutex
//
// Returns an error if rotation or writing fails.
func (s *RotatingFileSink) Write(data []byte) error {
	// Lock to ensure only one goroutine writes/rotates at a time
	s.mu.Lock()
	defer s.mu.Unlock() // Unlock when function returns (even if there's an error)

	// Check if adding this log entry would exceed the maximum size
	// len(data) gives the number of bytes in the data
	if s.size+int64(len(data)) > s.maxSizeMB {
		// File is too large, rotate it before writing
		if err := s.rotate(); err != nil {
			return err
		}
	}

	// Write the log data to the file
	// fmt.Fprintln adds a newline at the end
	n, err := fmt.Fprintln(s.file, string(data))
	if err != nil {
		return err
	}

	// Update our size counter with the number of bytes written
	s.size += int64(n)
	return nil
}

// rotate performs the file rotation process.
// This is called internally by Write when the file size limit is reached.
//
// Rotation steps:
//  1. Close the current file
//  2. Shift existing backups (app.log.1 → app.log.2, etc.)
//  3. Delete the oldest backup if it exceeds maxBackups
//  4. Rename current file to app.log.1
//  5. Create a new empty file
//
// Note: This method assumes the caller holds the mutex lock.
func (s *RotatingFileSink) rotate() error {
	// Close the current file before renaming it
	if err := s.file.Close(); err != nil {
		return err
	}

	// Shift existing backup files
	// Start from the highest number and work backwards
	// Example: if maxBackups=5, this shifts .4→.5, .3→.4, .2→.3, .1→.2
	for i := s.maxBackups - 1; i >= 1; i-- {
		oldPath := fmt.Sprintf("%s.%d", s.path, i)     // e.g., "app.log.3"
		newPath := fmt.Sprintf("%s.%d", s.path, i+1)   // e.g., "app.log.4"

		// Check if the old backup file exists
		if _, err := os.Stat(oldPath); err == nil {
			if i+1 > s.maxBackups {
				// This backup would exceed maxBackups, delete it
				os.Remove(oldPath)
			} else {
				// Rename the backup to the next number
				os.Rename(oldPath, newPath)
			}
		}
	}

	// Rename the current log file to .1 (most recent backup)
	backupPath := fmt.Sprintf("%s.1", s.path) // e.g., "app.log.1"
	if err := os.Rename(s.path, backupPath); err != nil {
		return err
	}

	// Create a new empty log file
	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Update the sink to use the new file
	s.file = file
	s.size = 0 // Reset size counter
	return nil
}

// Close closes the current log file and releases the file descriptor.
// This is thread-safe and can be called even if Write is in progress.
//
// Always call logger.Close() to ensure the file is properly closed.
func (s *RotatingFileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
