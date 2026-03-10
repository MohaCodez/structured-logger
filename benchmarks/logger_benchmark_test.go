package benchmarks

import (
	"testing"

	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
)

// NoOpSink discards all logs (for benchmarking)
type noOpSink struct{}

func (n *noOpSink) Write(data []byte) error {
	return nil
}

func (n *noOpSink) Close() error {
	return nil
}

// BenchmarkSyncLogging measures synchronous logging performance
func BenchmarkSyncLogging(b *testing.B) {
	config := logger.DefaultConfig()
	config.Level = logger.INFO
	config.Formatter = formatter.NewJSONFormatter()
	config.Sinks = []logger.Sink{&noOpSink{}}
	log := logger.NewWithConfig(config)
	defer log.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("benchmark_message")
	}
}

// BenchmarkAsyncLogging measures asynchronous logging performance
func BenchmarkAsyncLogging(b *testing.B) {
	config := logger.DefaultConfig()
	config.Level = logger.INFO
	config.Formatter = formatter.NewJSONFormatter()
	config.Sinks = []logger.Sink{&noOpSink{}}
	config.Async = true
	config.BufferSize = 1000
	log := logger.NewWithConfig(config)
	defer log.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("benchmark_message")
	}
}

// BenchmarkStructuredFields measures overhead of structured fields
func BenchmarkStructuredFields(b *testing.B) {
	config := logger.DefaultConfig()
	config.Level = logger.INFO
	config.Formatter = formatter.NewJSONFormatter()
	config.Sinks = []logger.Sink{&noOpSink{}}
	log := logger.NewWithConfig(config)
	defer log.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("benchmark_message",
			"user_id", 12345,
			"request_id", "abc123",
			"duration_ms", 45,
		)
	}
}

// BenchmarkJSONFormatting measures JSON formatting overhead
func BenchmarkJSONFormatting(b *testing.B) {
	formatter := formatter.NewJSONFormatter()
	entry := &logger.Entry{
		Timestamp: "2026-03-10T03:40:00+05:30",
		Level:     "INFO",
		Message:   "benchmark_message",
		Fields: map[string]interface{}{
			"user_id":    12345,
			"request_id": "abc123",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatter.Format(entry)
	}
}

// BenchmarkContextualLogging measures overhead of contextual logging
func BenchmarkContextualLogging(b *testing.B) {
	config := logger.DefaultConfig()
	config.Level = logger.INFO
	config.Formatter = formatter.NewJSONFormatter()
	config.Sinks = []logger.Sink{&noOpSink{}}
	baseLog := logger.NewWithConfig(config)
	defer baseLog.Close()

	// Create child logger with context
	childLog := baseLog.With("service", "api", "environment", "production")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		childLog.Info("benchmark_message")
	}
}

// BenchmarkWithCaller measures overhead of caller tracing
func BenchmarkWithCaller(b *testing.B) {
	config := logger.DefaultConfig()
	config.Level = logger.INFO
	config.Formatter = formatter.NewJSONFormatter()
	config.Sinks = []logger.Sink{&noOpSink{}}
	config.EnableCaller = true
	log := logger.NewWithConfig(config)
	defer log.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("benchmark_message")
	}
}

// BenchmarkLevelFiltering measures level filtering performance
func BenchmarkLevelFiltering(b *testing.B) {
	config := logger.DefaultConfig()
	config.Level = logger.ERROR // Set high threshold
	config.Formatter = formatter.NewJSONFormatter()
	config.Sinks = []logger.Sink{&noOpSink{}}
	log := logger.NewWithConfig(config)
	defer log.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Debug("filtered_message") // Should be filtered out
	}
}

// BenchmarkNestedContextualLogging measures nested child logger overhead
func BenchmarkNestedContextualLogging(b *testing.B) {
	config := logger.DefaultConfig()
	config.Level = logger.INFO
	config.Formatter = formatter.NewJSONFormatter()
	config.Sinks = []logger.Sink{&noOpSink{}}
	baseLog := logger.NewWithConfig(config)
	defer baseLog.Close()

	serviceLog := baseLog.With("service", "api")
	requestLog := serviceLog.With("request_id", "abc123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		requestLog.Info("benchmark_message")
	}
}
