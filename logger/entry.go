package logger

import (
	"encoding/json"
	"time"
)

type Entry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

func newEntry(level Level, message string) *Entry {
	return &Entry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level.String(),
		Message:   message,
	}
}

func (e *Entry) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
