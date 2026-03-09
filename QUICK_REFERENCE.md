# Quick Reference

## Installation
```bash
go get github.com/MohaCodez/structured-logger
```

## Basic Usage
```go
import "github.com/MohaCodez/structured-logger/logger"

log := logger.New(logger.INFO)
defer log.Close()

log.Info("message", "key", "value")
```

## Log Levels
```go
log.Debug("debug")   // 0 - Most verbose
log.Info("info")     // 1 - General info
log.Warn("warning")  // 2 - Warnings
log.Error("error")   // 3 - Errors
log.Fatal("fatal")   // 4 - Critical (exits)
```

## Configuration
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
```

## Default Config
```go
config := logger.DefaultConfig()
config.Level = logger.DEBUG
config.EnableCaller = true
log := logger.NewWithConfig(config)
```

## Multiple Sinks
```go
consoleSink := sink.NewConsoleSink()
fileSink, _ := sink.NewFileSink("app.log")

config := logger.Config{
    Sinks: []logger.Sink{consoleSink, fileSink},
}
log := logger.NewWithConfig(config)
defer log.Close()  // Closes all sinks
```

## Async Mode
```go
config := logger.Config{
    Async:      true,
    BufferSize: 500,
}
log := logger.NewWithConfig(config)
defer log.Close()  // Flushes queue
```

## Structured Fields
```go
log.Info("user_login",
    "user_id", 123,
    "username", "alice",
    "ip", "10.1.2.4",
    "success", true,
)
```

## Caller Tracing
```go
config := logger.DefaultConfig()
config.EnableCaller = true
log := logger.NewWithConfig(config)

log.Error("error", "msg", "failed")
// Output includes: "caller": "main.go:42"
```

## Custom Formatter
```go
type MyFormatter struct{}

func (f *MyFormatter) Format(entry *logger.Entry) ([]byte, error) {
    return []byte(fmt.Sprintf("[%s] %s\n", entry.Level, entry.Message)), nil
}

config := logger.Config{
    Formatter: &MyFormatter{},
}
```

## Custom Sink
```go
type MySink struct{}

func (s *MySink) Write(data []byte) error {
    // Custom output logic
    return nil
}

func (s *MySink) Close() error {
    return nil
}

config := logger.Config{
    Sinks: []logger.Sink{&MySink{}},
}
```

## Testing
```go
type testSink struct {
    logs []string
}

func (t *testSink) Write(data []byte) error {
    t.logs = append(t.logs, string(data))
    return nil
}

func (t *testSink) Close() error {
    return nil
}

// In test
sink := &testSink{}
log := logger.NewWithSinks(logger.DEBUG, formatter, []logger.Sink{sink})
```

## Constructors

| Constructor | Parameters | Use Case |
|-------------|------------|----------|
| `New(level)` | level | Simple, defaults |
| `NewWithConfig(config)` | config | Full control |
| `NewWithFormatter(level, formatter)` | level, formatter | Custom format |
| `NewWithSinks(level, formatter, sinks)` | level, formatter, sinks | Multiple outputs |
| `NewWithCaller(level, formatter, sinks, caller)` | level, formatter, sinks, caller | With tracing |
| `NewAsync(level, formatter, sinks, caller, bufSize)` | level, formatter, sinks, caller, bufSize | Async mode |

## Config Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| Level | Level | INFO | Minimum log level |
| Formatter | Formatter | defaultFormatter | Output formatter |
| Sinks | []Sink | [consoleSink] | Output destinations |
| EnableCaller | bool | false | Include file:line |
| Async | bool | false | Non-blocking mode |
| BufferSize | int | 100 | Async queue size |

## Output Format (JSON)
```json
{
  "timestamp": "2026-03-10T03:14:00+05:30",
  "level": "INFO",
  "message": "user_login",
  "caller": "main.go:42",
  "user_id": 123,
  "ip": "10.1.2.4"
}
```

## Best Practices

✓ Always `defer log.Close()`  
✓ Use structured fields, not string formatting  
✓ Set INFO or WARN level in production  
✓ Enable async for high throughput  
✓ Disable caller in production  
✓ Include request IDs for tracing  
✓ Add timing information  

## Performance

| Mode | Throughput | Use Case |
|------|------------|----------|
| Sync | 68K logs/sec | Low volume |
| Async | 100K logs/sec | High volume |

## Common Patterns

**Request ID:**
```go
log.Info("request", "request_id", reqID)
```

**Timing:**
```go
start := time.Now()
// ... operation
log.Info("done", "duration_ms", time.Since(start).Milliseconds())
```

**Error Context:**
```go
log.Error("failed", "error", err.Error(), "user_id", userID)
```

## Running Tests
```bash
go test ./logger ./formatter ./sink ./async -v
go test ./logger ./formatter ./sink ./async -cover
```

## Documentation
- `README.md` - Full documentation
- `ARCHITECTURE.md` - Design details
- `INTEGRATION.md` - Integration guide
- `examples/` - Working examples
