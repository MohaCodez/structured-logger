package sink

import (
	"fmt"
	"os"
)

type FileSink struct {
	file *os.File
}

func NewFileSink(path string) (*FileSink, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return &FileSink{file: file}, nil
}

func (s *FileSink) Write(data []byte) error {
	_, err := fmt.Fprintln(s.file, string(data))
	return err
}

func (s *FileSink) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
