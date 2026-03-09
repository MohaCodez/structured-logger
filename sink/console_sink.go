package sink

import (
	"fmt"
	"os"
)

type ConsoleSink struct{}

func NewConsoleSink() *ConsoleSink {
	return &ConsoleSink{}
}

func (s *ConsoleSink) Write(data []byte) error {
	_, err := fmt.Fprintln(os.Stdout, string(data))
	return err
}

func (s *ConsoleSink) Close() error {
	return nil
}
