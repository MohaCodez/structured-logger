package sink

import (
	"fmt"
	"os"
)

// ConsoleSink writes log output to the console (standard output).
// This is the default sink and is useful for:
//   - Development and debugging (see logs immediately)
//   - Containerized applications (logs go to stdout, captured by container runtime)
//   - Simple applications that don't need persistent logs
//
// ConsoleSink has no state and is safe for concurrent use.
type ConsoleSink struct{}

// NewConsoleSink creates a new console sink.
// This is the standard way to create a ConsoleSink instance.
//
// Example:
//
//	config := logger.DefaultConfig()
//	config.Sinks = []logger.Sink{sink.NewConsoleSink()}
//	log := logger.NewWithConfig(config)
func NewConsoleSink() *ConsoleSink {
	return &ConsoleSink{}
}

// Write outputs the log data to stdout (standard output / console).
// Each log entry is written as a single line.
//
// Returns an error if writing to stdout fails (rare, but possible if
// stdout is redirected to a full disk or closed pipe).
func (s *ConsoleSink) Write(data []byte) error {
	// fmt.Fprintln writes to os.Stdout and adds a newline
	// os.Stdout is the standard output stream (usually the terminal)
	_, err := fmt.Fprintln(os.Stdout, string(data))
	return err
}

// Close does nothing for console output.
// stdout is managed by the operating system and doesn't need to be closed.
// This method exists to satisfy the Sink interface.
func (s *ConsoleSink) Close() error {
	return nil
}
