# Structured Logger (Go)

A modular, plug-and-play structured logging library in Go featuring leveled logging, JSON output, caller tracing, and extensible output sinks.

[![Tests](https://img.shields.io/badge/tests-41%20passing-brightgreen)]()
[![Coverage](https://img.shields.io/badge/coverage-85%25-green)]()
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue)]()

---

## Features

✓ **Structured JSON logging** - Machine-readable log events  
✓ **Log levels** - DEBUG, INFO, WARN, ERROR, FATAL  
✓ **Contextual logging** - Child loggers with inherited fields  
✓ **Caller tracing** - Optional file:line information  
✓ **Pluggable sinks** - Console, file, or custom outputs  
✓ **File rotation** - Size-based log rotation  
✓ **Async logging** - Non-blocking high-throughput mode  
✓ **Flexible formatters** - JSON or custom formats  
✓ **Minimal dependencies** - Pure Go implementation  
✓ **Production-ready** - Comprehensive test coverage  

---

## Quick Start

### Installation

```bash
go get github.com/MohaCodez/structured-logger
```

### Basic Usage

```go
package main

import "github.com/MohaCodez/structured-logger/logger"

func main() {
    log := logger.New(logger.INFO)
    defer log.Close()

    log.Info("user_login",
        "user_id", 123,
        "ip", "10.1.2.4",
    )
}
```

**Output:**
```json
{"timestamp":"2026-03-10T03:14:00+05:30","level":"INFO","message":"user_login","user_id":123,"ip":"10.1.2.4"}
```

---

## Architecture

```
Application
     │
     ▼
Logger API (Debug/Info/Warn/Error/Fatal)
     │
     ▼
Entry Builder (timestamp, level, message, fields, caller)
     │
     ▼
Level Filter (skip if below threshold)
     │
     ▼
Formatter (JSON, custom)
     │
     ▼
Async Worker (optional)
     │
     ▼
Sink Dispatcher (fan-out to multiple sinks)
     │
     ├── Console Sink
     ├── File Sink
     └── Custom Sink
```

---

## Configuration

### Using Config Struct (Recommended)

```go
import (
    "github.com/MohaCodez/structured-logger/formatter"
    "github.com/MohaCodez/structured-logger/logger"
    "github.com/MohaCodez/structured-logger/sink"
)

config := logger.Config{
    Level:        logger.INFO,
    Formatter:    formatter.NewJSONFormatter(),
    Sinks:        []logger.Sink{sink.NewConsoleSink()},
    EnableCaller: true,
    Async:        false,
    BufferSize:   100,
}

log := logger.NewWithConfig(config)
defer log.Close()
```

### Default Configuration

```go
config := logger.DefaultConfig()
config.Level = logger.DEBUG
config.EnableCaller = true

log := logger.NewWithConfig(config)
```

---

## Log Levels

```go
log.Debug("debug message")   // Development details
log.Info("info message")      // General information
log.Warn("warning message")   // Warning conditions
log.Error("error message")    // Error conditions
log.Fatal("fatal message")    // Critical errors (exits program)
```

**Level Filtering:**
```go
log := logger.New(logger.WARN)  // Only WARN, ERROR, FATAL will be logged
log.Info("ignored")             // Won't be logged
log.Error("logged")             // Will be logged
```

---

## Structured Fields

Add context to logs with key/value pairs:

```go
log.Info("payment_processed",
    "transaction_id", "txn_12345",
    "amount", 99.99,
    "currency", "USD",
    "user_id", 123,
    "success", true,
)
```

**Output:**
```json
{
  "timestamp": "2026-03-10T03:14:00+05:30",
  "level": "INFO",
  "message": "payment_processed",
  "transaction_id": "txn_12345",
  "amount": 99.99,
  "currency": "USD",
  "user_id": 123,
  "success": true
}
```

---

## Contextual Logging

Create child loggers with inherited fields for consistent metadata:

```go
// Base logger
baseLog := logger.New(logger.INFO)

// Service-level logger with context
serviceLog := baseLog.With("service", "auth", "environment", "production")

// Request-level logger inherits service context
requestLog := serviceLog.With("request_id", "abc123", "user_id", 42)

requestLog.Info("processing_request")
// Output includes: service, environment, request_id, user_id
```

**Output:**
```json
{
  "timestamp": "2026-03-10T03:14:00+05:30",
  "level": "INFO",
  "message": "processing_request",
  "service": "auth",
  "environment": "production",
  "request_id": "abc123",
  "user_id": 42
}
```

**Benefits:**
- No need to repeat fields on every log call
- Parent logger remains unchanged
- Nested contexts work seamlessly
- Call fields override context fields

---

## Caller Tracing

Enable file and line number tracking:

```go
config := logger.DefaultConfig()
config.EnableCaller = true
log := logger.NewWithConfig(config)

log.Error("database_error", "error", "connection timeout")
```

**Output:**
```json
{
  "timestamp": "2026-03-10T03:14:00+05:30",
  "level": "ERROR",
  "message": "database_error",
  "caller": "main.go:42",
  "error": "connection timeout"
}
```

---

## Multiple Sinks

Write logs to multiple destinations simultaneously:

```go
consoleSink := sink.NewConsoleSink()
fileSink, _ := sink.NewFileSink("app.log")

config := logger.Config{
    Level:     logger.INFO,
    Formatter: formatter.NewJSONFormatter(),
    Sinks:     []logger.Sink{consoleSink, fileSink},
}

log := logger.NewWithConfig(config)
defer log.Close()  // Closes all sinks
```

---

## File Rotation

Prevent log files from growing indefinitely with size-based rotation:

```go
import "github.com/MohaCodez/structured-logger/sink"

// Create rotating file sink
// MaxSize: 10 MB, MaxBackups: 5
rotatingSink, err := sink.NewRotatingFileSink("app.log", 10, 5)
if err != nil {
    panic(err)
}

config := logger.Config{
    Level:     logger.INFO,
    Formatter: formatter.NewJSONFormatter(),
    Sinks:     []logger.Sink{rotatingSink},
}

log := logger.NewWithConfig(config)
defer log.Close()
```

**Rotation behavior:**
- When `app.log` exceeds 10 MB, it's renamed to `app.log.1`
- Previous backups shift: `app.log.1` → `app.log.2`, etc.
- Keeps maximum of 5 backup files
- Oldest backups are automatically deleted

**Files created:**
```
app.log       (current log file)
app.log.1     (most recent backup)
app.log.2
app.log.3
app.log.4
app.log.5     (oldest backup)
```

---

## Asynchronous Logging

Enable non-blocking logging for high-throughput systems:

```go
config := logger.Config{
    Level:      logger.INFO,
    Formatter:  formatter.NewJSONFormatter(),
    Sinks:      []logger.Sink{sink.NewConsoleSink()},
    Async:      true,
    BufferSize: 500,  // Queue size
}

log := logger.NewWithConfig(config)
defer log.Close()  // Flushes queue before closing

// Non-blocking log calls
for i := 0; i < 10000; i++ {
    log.Info("high_throughput", "iteration", i)
}
```

**Performance:** ~1.5x faster than synchronous logging

---

## Custom Formatters

Implement the `Formatter` interface:

```go
type Formatter interface {
    Format(entry *logger.Entry) ([]byte, error)
}
```

Example: Text formatter
```go
type TextFormatter struct{}

func (f *TextFormatter) Format(entry *logger.Entry) ([]byte, error) {
    return []byte(fmt.Sprintf("[%s] %s: %s\n", 
        entry.Level, entry.Timestamp, entry.Message)), nil
}
```

---

## Custom Sinks

Implement the `Sink` interface:

```go
type Sink interface {
    Write(data []byte) error
    Close() error
}
```

Example: HTTP sink
```go
type HTTPSink struct {
    url string
}

func (s *HTTPSink) Write(data []byte) error {
    _, err := http.Post(s.url, "application/json", bytes.NewReader(data))
    return err
}

func (s *HTTPSink) Close() error {
    return nil
}
```

---

## Complete Example

```go
package main

import (
    "github.com/MohaCodez/structured-logger/formatter"
    "github.com/MohaCodez/structured-logger/logger"
    "github.com/MohaCodez/structured-logger/sink"
)

func main() {
    // Setup
    fileSink, err := sink.NewFileSink("app.log")
    if err != nil {
        panic(err)
    }

    config := logger.Config{
        Level:        logger.DEBUG,
        Formatter:    formatter.NewJSONFormatter(),
        Sinks:        []logger.Sink{sink.NewConsoleSink(), fileSink},
        EnableCaller: true,
        Async:        true,
        BufferSize:   200,
    }

    log := logger.NewWithConfig(config)
    defer log.Close()

    // Application logs
    log.Info("server_started", "port", 8080)
    
    log.Debug("processing_request",
        "method", "GET",
        "path", "/api/users",
        "duration_ms", 45,
    )

    log.Warn("rate_limit_exceeded",
        "user_id", 123,
        "limit", 100,
        "current", 150,
    )

    log.Error("database_error",
        "operation", "SELECT",
        "table", "users",
        "error", "connection timeout",
    )
}
```

---

## Project Structure

```
structured-logger/
├── logger/           # Core logger implementation
│   ├── logger.go
│   ├── level.go
│   ├── entry.go
│   └── config.go
├── formatter/        # Output formatters
│   ├── formatter.go
│   └── json_formatter.go
├── sink/            # Output destinations
│   ├── sink.go
│   ├── console_sink.go
│   ├── file_sink.go
│   └── rotating_file_sink.go
├── async/           # Async worker
│   └── worker.go
├── benchmarks/      # Performance benchmarks
│   └── logger_benchmark_test.go
├── examples/        # Usage examples
└── README.md
```

---

## Testing

Run tests:
```bash
go test ./logger ./formatter ./sink ./async -v
```

With coverage:
```bash
go test ./logger ./formatter ./sink ./async -cover
```

**Test Results:** 41/41 tests passing, 85% average coverage

---

## Performance Benchmarks

Benchmark results on Intel Core i7-9750H @ 2.60GHz, Ubuntu 22.04 LTS, Go 1.25.7:

| Operation | Time/op | Throughput | Memory/op | Allocs/op |
|-----------|---------|------------|-----------|-----------|
| Sync Logging | 2,144 ns | 466K logs/sec | 912 B | 17 |
| Async Logging | 2,010 ns | 498K logs/sec | 928 B | 18 |
| Structured Fields | 4,141 ns | 241K logs/sec | 1,745 B | 25 |
| Contextual Logging | 3,292 ns | 304K logs/sec | 1,392 B | 22 |
| Level Filtering | 1.9 ns | 526M ops/sec | 0 B | 0 |

**Key Insights:**
- Level filtering is essentially free (zero allocations)
- Async mode provides ~6% better performance in high-throughput scenarios
- Structured fields add ~93% overhead due to field processing
- Contextual logging adds ~54% overhead for field inheritance

The async advantage becomes more pronounced under I/O pressure, where non-blocking writes prevent caller delays.

Run benchmarks:
```bash
go test ./benchmarks -bench=. -benchmem
```

See [BENCHMARKS.md](BENCHMARKS.md) for detailed analysis.

---

## API Reference

### Logger Methods

| Method | Description |
|--------|-------------|
| `Debug(msg, ...fields)` | Log debug message |
| `Info(msg, ...fields)` | Log info message |
| `Warn(msg, ...fields)` | Log warning message |
| `Error(msg, ...fields)` | Log error message |
| `Fatal(msg, ...fields)` | Log fatal and exit |
| `With(...fields)` | Create child logger with context |
| `Close()` | Close logger and sinks |

### Constructors

| Constructor | Use Case |
|-------------|----------|
| `New(level)` | Simple logger with defaults |
| `NewWithConfig(config)` | Full configuration control |

---

## Best Practices

1. **Always defer Close()** - Ensures logs are flushed
   ```go
   log := logger.NewWithConfig(config)
   defer log.Close()
   ```

2. **Use structured fields** - Better than string formatting
   ```go
   // Good
   log.Info("user_login", "user_id", 123)
   
   // Avoid
   log.Info(fmt.Sprintf("user %d logged in", 123))
   ```

3. **Set appropriate log levels** - Use INFO or WARN in production
   ```go
   config := logger.DefaultConfig()
   config.Level = logger.INFO  // Production
   ```

4. **Enable async for high throughput** - Reduces I/O blocking
   ```go
   config.Async = true
   config.BufferSize = 500
   ```

5. **Use caller tracing in development** - Disable in production for performance
   ```go
   config.EnableCaller = (env == "development")
   ```

---

## Context Integration

Store and retrieve loggers from context.Context for request-scoped logging:

```go
import (
    "context"
    "github.com/MohaCodez/structured-logger/logger"
)

// Store logger in context
ctx = logger.WithContext(ctx, requestLog)

// Retrieve logger from context
log := logger.FromContext(ctx)
log.Info("handling_request")
```

**Usage Pattern:**
```go
func handleRequest(ctx context.Context) {
    log := logger.FromContext(ctx)
    log.Info("processing_request")
    
    // Pass context to other functions
    authenticateUser(ctx)
}

func authenticateUser(ctx context.Context) {
    log := logger.FromContext(ctx)
    log.Info("authenticating_user")
}
```

---

## Async Buffer Policies

Control behavior when async buffer is full:

```go
config := logger.DefaultConfig()
config.Async = true
config.BufferSize = 100
config.BufferFullPolicy = logger.BlockOnFull  // Default: blocks caller
// config.BufferFullPolicy = logger.DropOnFull  // Alternative: drops logs

log := logger.NewWithConfig(config)
```

**Policies:**
- `BlockOnFull` (default): Provides backpressure, ensures no log loss
- `DropOnFull`: Non-blocking, may drop logs under extreme load

---

## Roadmap

- [x] Core logging engine
- [x] Structured fields
- [x] Formatter abstraction
- [x] Pluggable sinks
- [x] Caller tracing
- [x] Async logging
- [x] Configuration system
- [x] Comprehensive tests
- [x] Contextual logging (child loggers)
- [x] Log rotation
- [x] Performance benchmarks
- [ ] Log sampling
- [ ] Distributed tracing integration
- [ ] Cloud logging sinks (AWS CloudWatch, GCP Logging)

---

## License

MIT License

---

## Contributing

Contributions welcome! Please ensure:
- All tests pass
- Code coverage remains above 80%
- Follow existing code style
- Add tests for new features

---

## Support

- Issues: [GitHub Issues](https://github.com/MohaCodez/structured-logger/issues)
- Documentation: See `examples/` directory
- Tests: See `*_test.go` files
