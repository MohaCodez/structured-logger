package logger

import (
	"fmt"
	"os"
)

type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

type Sink interface {
	Write(data []byte) error
	Close() error
}

type Logger struct {
	level        Level
	formatter    Formatter
	sinks        []Sink
	enableCaller bool
}

func New(level Level) *Logger {
	return &Logger{
		level:        level,
		formatter:    &defaultFormatter{},
		sinks:        []Sink{&defaultConsoleSink{}},
		enableCaller: false,
	}
}

func NewWithFormatter(level Level, formatter Formatter) *Logger {
	return &Logger{
		level:        level,
		formatter:    formatter,
		sinks:        []Sink{&defaultConsoleSink{}},
		enableCaller: false,
	}
}

func NewWithSinks(level Level, formatter Formatter, sinks []Sink) *Logger {
	return &Logger{
		level:        level,
		formatter:    formatter,
		sinks:        sinks,
		enableCaller: false,
	}
}

func NewWithCaller(level Level, formatter Formatter, sinks []Sink, enableCaller bool) *Logger {
	return &Logger{
		level:        level,
		formatter:    formatter,
		sinks:        sinks,
		enableCaller: enableCaller,
	}
}

func (l *Logger) Close() error {
	for _, sink := range l.sinks {
		if err := sink.Close(); err != nil {
			return err
		}
	}
	return nil
}

type defaultFormatter struct{}

func (f *defaultFormatter) Format(entry *Entry) ([]byte, error) {
	return []byte(fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s"}`, 
		entry.Timestamp, entry.Level, entry.Message)), nil
}

type defaultConsoleSink struct{}

func (s *defaultConsoleSink) Write(data []byte) error {
	fmt.Println(string(data))
	return nil
}

func (s *defaultConsoleSink) Close() error {
	return nil
}

func (l *Logger) log(level Level, message string, keyValues ...interface{}) {
	if level < l.level {
		return
	}

	fields := parseFields(keyValues)
	entry := newEntry(level, message, fields, l.enableCaller)
	data, err := l.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to format log entry: %v\n", err)
		return
	}

	for _, sink := range l.sinks {
		if err := sink.Write(data); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to sink: %v\n", err)
		}
	}
}

func parseFields(keyValues []interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	
	if len(keyValues)%2 != 0 {
		fmt.Fprintf(os.Stderr, "warning: uneven number of key/value pairs, ignoring last element\n")
		keyValues = keyValues[:len(keyValues)-1]
	}

	for i := 0; i < len(keyValues); i += 2 {
		key, ok := keyValues[i].(string)
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: non-string key at position %d, skipping pair\n", i)
			continue
		}
		fields[key] = keyValues[i+1]
	}

	return fields
}

func (l *Logger) Debug(message string, keyValues ...interface{}) {
	l.log(DEBUG, message, keyValues...)
}

func (l *Logger) Info(message string, keyValues ...interface{}) {
	l.log(INFO, message, keyValues...)
}

func (l *Logger) Warn(message string, keyValues ...interface{}) {
	l.log(WARN, message, keyValues...)
}

func (l *Logger) Error(message string, keyValues ...interface{}) {
	l.log(ERROR, message, keyValues...)
}

func (l *Logger) Fatal(message string, keyValues ...interface{}) {
	l.log(FATAL, message, keyValues...)
	os.Exit(1)
}
