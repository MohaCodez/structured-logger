package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	level Level
}

func New(level Level) *Logger {
	return &Logger{level: level}
}

func (l *Logger) log(level Level, message string) {
	if level < l.level {
		return
	}

	entry := newEntry(level, message)
	data, err := entry.ToJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal log entry: %v\n", err)
		return
	}

	fmt.Println(string(data))
}

func (l *Logger) Debug(message string) {
	l.log(DEBUG, message)
}

func (l *Logger) Info(message string) {
	l.log(INFO, message)
}

func (l *Logger) Warn(message string) {
	l.log(WARN, message)
}

func (l *Logger) Error(message string) {
	l.log(ERROR, message)
}

func (l *Logger) Fatal(message string) {
	l.log(FATAL, message)
	os.Exit(1)
}
