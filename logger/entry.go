package logger

import "time"

type Entry struct {
	Timestamp string
	Level     string
	Message   string
	Fields    map[string]interface{}
}

func newEntry(level Level, message string, fields map[string]interface{}) *Entry {
	return &Entry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level.String(),
		Message:   message,
		Fields:    fields,
	}
}
