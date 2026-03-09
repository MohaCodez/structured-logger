package logger

import "github.com/MohaCodez/structured-logger/async"

type Config struct {
	Level        Level
	Formatter    Formatter
	Sinks        []Sink
	EnableCaller bool
	Async        bool
	BufferSize   int
}

func DefaultConfig() Config {
	return Config{
		Level:        INFO,
		Formatter:    &defaultFormatter{},
		Sinks:        []Sink{&defaultConsoleSink{}},
		EnableCaller: false,
		Async:        false,
		BufferSize:   100,
	}
}

func NewWithConfig(config Config) *Logger {
	logger := &Logger{
		level:        config.Level,
		formatter:    config.Formatter,
		sinks:        config.Sinks,
		enableCaller: config.EnableCaller,
		asyncWorker:  nil,
	}

	if config.Async {
		logger.asyncWorker = async.NewWorker(config.BufferSize)
	}

	return logger
}
