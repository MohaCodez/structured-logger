package logger

import (
	"fmt"
	"os"

	"github.com/MohaCodez/structured-logger/async"
)

type Formatter interface {
	Format(entry *Entry) ([]byte, error)
}

type Sink interface {
	Write(data []byte) error
	Close() error
}

type Logger struct {
	level            Level
	formatter        Formatter
	sinks            []Sink
	enableCaller     bool
	asyncWorker      *async.Worker
	contextFields    map[string]interface{}
	exitFunc         func(int)
	sinkErrorHandler func(error)
}

func New(level Level) *Logger {
	config := DefaultConfig()
	config.Level = level
	return NewWithConfig(config)
}

func (l *Logger) Close() error {
	if l.asyncWorker != nil {
		l.asyncWorker.Stop()
	}
	for _, sink := range l.sinks {
		if err := sink.Close(); err != nil {
			return err
		}
	}
	return nil
}

// With creates a child logger with additional context fields.
// The parent logger remains unchanged.
func (l *Logger) With(keyValues ...interface{}) *Logger {
	fields := parseFields(keyValues)
	
	// Create new context fields map with parent fields + new fields
	newContextFields := make(map[string]interface{}, len(l.contextFields)+len(fields))
	for k, v := range l.contextFields {
		newContextFields[k] = v
	}
	for k, v := range fields {
		newContextFields[k] = v
	}
	
	// Return new logger with merged context
	return &Logger{
		level:            l.level,
		formatter:        l.formatter,
		sinks:            l.sinks,
		enableCaller:     l.enableCaller,
		asyncWorker:      l.asyncWorker,
		contextFields:    newContextFields,
		exitFunc:         l.exitFunc,
		sinkErrorHandler: l.sinkErrorHandler,
	}
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
	
	// Merge context fields with call fields (call fields override)
	mergedFields := make(map[string]interface{}, len(l.contextFields)+len(fields))
	for k, v := range l.contextFields {
		mergedFields[k] = v
	}
	for k, v := range fields {
		mergedFields[k] = v
	}
	
	entry := newEntry(level, message, mergedFields, l.enableCaller)
	data, err := l.formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to format log entry: %v\n", err)
		return
	}

	if l.asyncWorker != nil {
		// Async mode: enqueue for background processing
		sinksCopy := make([]async.Sink, len(l.sinks))
		for i, s := range l.sinks {
			sinksCopy[i] = s
		}
		l.asyncWorker.Enqueue(data, sinksCopy)
	} else {
		// Sync mode: write immediately
		for _, sink := range l.sinks {
			if err := sink.Write(data); err != nil {
				l.sinkErrorHandler(err)
			}
		}
	}
}

func parseFields(keyValues []interface{}) map[string]interface{} {
	fields := make(map[string]interface{})
	
	if len(keyValues)%2 != 0 {
		fmt.Fprintf(os.Stderr, "structured-logger: odd number of fields passed to log call, last key has no value\n")
		keyValues = append(keyValues, "MISSING_VALUE")
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
	l.Close()
	l.exitFunc(1)
}
