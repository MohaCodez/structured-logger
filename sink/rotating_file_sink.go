package sink

import (
	"fmt"
	"os"
	"sync"
)

type RotatingFileSink struct {
	path       string
	maxSizeMB  int64
	maxBackups int
	file       *os.File
	size       int64
	mu         sync.Mutex
}

func NewRotatingFileSink(path string, maxSizeMB int, maxBackups int) (*RotatingFileSink, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Get current file size
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to stat log file: %w", err)
	}

	return &RotatingFileSink{
		path:       path,
		maxSizeMB:  int64(maxSizeMB) * 1024 * 1024,
		maxBackups: maxBackups,
		file:       file,
		size:       info.Size(),
	}, nil
}

func (s *RotatingFileSink) Write(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if rotation needed
	if s.size+int64(len(data)) > s.maxSizeMB {
		if err := s.rotate(); err != nil {
			return err
		}
	}

	// Write to file
	n, err := fmt.Fprintln(s.file, string(data))
	if err != nil {
		return err
	}

	s.size += int64(n)
	return nil
}

func (s *RotatingFileSink) rotate() error {
	// Close current file
	if err := s.file.Close(); err != nil {
		return err
	}

	// Shift existing backups
	for i := s.maxBackups - 1; i >= 1; i-- {
		oldPath := fmt.Sprintf("%s.%d", s.path, i)
		newPath := fmt.Sprintf("%s.%d", s.path, i+1)

		if _, err := os.Stat(oldPath); err == nil {
			if i+1 > s.maxBackups {
				// Delete oldest backup
				os.Remove(oldPath)
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}

	// Rename current file to .1
	backupPath := fmt.Sprintf("%s.1", s.path)
	if err := os.Rename(s.path, backupPath); err != nil {
		return err
	}

	// Create new file
	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	s.file = file
	s.size = 0
	return nil
}

func (s *RotatingFileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
