package logger

import (
	"fmt"
	"os"

	"github.com/MohaCodez/structured-logger/async"
)

// BufferFullPolicy defines behavior when async buffer is full
type BufferFullPolicy int

const (
	// BlockOnFull blocks the caller until buffer space is available (default)
	BlockOnFull BufferFullPolicy = iota
	// DropOnFull drops the log entry and continues without blocking
	DropOnFull
)

type Config struct {
	Level            Level
	Formatter        Formatter
	Sinks            []Sink
	EnableCaller     bool
	Async            bool
	BufferSize       int
	ExitFunc         func(int) // Function to call on Fatal, defaults to os.Exit
	BufferFullPolicy BufferFullPolicy
	SinkErrorHandler func(error) // Optional handler for sink write errors
}

func DefaultConfig() Config {
	return Config{
		Level:            INFO,
		Formatter:        &defaultFormatter{},
		Sinks:            []Sink{&defaultConsoleSink{}},
		EnableCaller:     false,
		Async:            false,
		BufferSize:       100,
		ExitFunc:         os.Exit,
		BufferFullPolicy: BlockOnFull,
	}
}

func NewWithConfig(config Config) *Logger {
	exitFunc := config.ExitFunc
	if exitFunc == nil {
		exitFunc = os.Exit
	}

	sinkErrorHandler := config.SinkErrorHandler
	if sinkErrorHandler == nil {
		sinkErrorHandler = func(err error) {
			fmt.Fprintf(os.Stderr, "sink write error: %v\n", err)
		}
	}

	logger := &Logger{
		level:            config.Level,
		formatter:        config.Formatter,
		sinks:            config.Sinks,
		enableCaller:     config.EnableCaller,
		asyncWorker:      nil,
		contextFields:    make(map[string]interface{}),
		exitFunc:         exitFunc,
		sinkErrorHandler: sinkErrorHandler,
	}

	if config.Async {
		logger.asyncWorker = async.NewWorker(config.BufferSize, config.BufferFullPolicy == DropOnFull, sinkErrorHandler)
	}

	return logger
}
