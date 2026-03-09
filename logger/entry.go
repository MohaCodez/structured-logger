package logger

import (
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

type Entry struct {
	Timestamp string
	Level     string
	Message   string
	Caller    string
	Fields    map[string]interface{}
}

func newEntry(level Level, message string, fields map[string]interface{}, enableCaller bool) *Entry {
	entry := &Entry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level.String(),
		Message:   message,
		Fields:    fields,
	}

	if enableCaller {
		entry.Caller = getCaller()
	}

	return entry
}

func getCaller() string {
	_, file, line, ok := runtime.Caller(4) // skip: getCaller, newEntry, log, Debug/Info/etc
	if !ok {
		return "unknown"
	}
	return filepath.Base(file) + ":" + strconv.Itoa(line)
}
