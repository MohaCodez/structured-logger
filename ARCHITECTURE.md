# Architecture Overview

## System Design

The structured logger follows a **pipeline architecture** where log events flow through multiple stages before reaching their destination.

```
┌─────────────┐
│ Application │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│         Logger API                  │
│  Debug/Info/Warn/Error/Fatal        │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│       Entry Builder                 │
│  - Timestamp                        │
│  - Level                            │
│  - Message                          │
│  - Fields (key/value)               │
│  - Caller (optional)                │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│       Level Filter                  │
│  Skip if level < threshold          │
└──────┬──────────────────────────────┘
       │
       ▼
┌─────────────────────────────────────┐
│       Formatter                     │
│  Convert Entry → []byte             │
│  (JSON, Text, Custom)               │
└──────┬──────────────────────────────┘
       │
       ├─── Sync Mode ───┐
       │                 │
       │                 ▼
       │         ┌───────────────┐
       │         │ Sink Dispatch │
       │         └───────┬───────┘
       │                 │
       └─ Async Mode ────┤
                         │
                         ▼
                 ┌───────────────┐
                 │ Async Worker  │
                 │ (Queue)       │
                 └───────┬───────┘
                         │
                         ▼
                 ┌───────────────┐
                 │ Sink Dispatch │
                 └───────┬───────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
    ┌────────┐      ┌────────┐     ┌────────┐
    │Console │      │  File  │     │ Custom │
    │  Sink  │      │  Sink  │     │  Sink  │
    └────────┘      └────────┘     └────────┘
```

---

## Component Responsibilities

### 1. Logger API
- Public interface for application code
- Methods: `Debug()`, `Info()`, `Warn()`, `Error()`, `Fatal()`
- Accepts message and variadic key/value pairs
- Manages lifecycle (initialization, close)

### 2. Entry Builder
- Creates structured log entry
- Captures timestamp (RFC3339)
- Converts level to string
- Parses key/value fields into map
- Optionally captures caller info (file:line)

### 3. Level Filter
- Compares entry level against logger threshold
- Skips processing if below minimum level
- Enables efficient filtering without formatting overhead

### 4. Formatter
- Converts Entry struct to byte array
- Default: JSON formatter
- Extensible via Formatter interface
- Handles field serialization

### 5. Async Worker (Optional)
- Buffered channel queue
- Background goroutine processes entries
- Non-blocking enqueue
- Graceful shutdown with flush

### 6. Sink Dispatcher
- Fan-out to multiple sinks
- Writes formatted data to each sink
- Error handling per sink (doesn't fail all on one error)

### 7. Sinks
- Final output destinations
- Console: stdout
- File: append mode with file handle
- Custom: implement Sink interface

---

## Data Flow

### Synchronous Mode
```
log.Info("msg", "key", "val")
  → Entry created
  → Level checked
  → Formatted to JSON
  → Written to all sinks immediately
  → Returns to caller
```

### Asynchronous Mode
```
log.Info("msg", "key", "val")
  → Entry created
  → Level checked
  → Formatted to JSON
  → Enqueued to channel
  → Returns to caller immediately
  
Background worker:
  → Dequeues entry
  → Writes to all sinks
```

---

## Key Design Decisions

### 1. Interface-Based Extensibility
- `Formatter` interface allows custom output formats
- `Sink` interface allows custom destinations
- Easy to add new implementations without modifying core

### 2. Separation of Concerns
- Entry creation separate from formatting
- Formatting separate from output
- Each component has single responsibility

### 3. Configuration Over Constructors
- `Config` struct provides unified configuration
- Easier to customize than multiple constructor parameters
- Backward compatible with simple constructors

### 4. Optional Features
- Caller tracing: disabled by default (performance)
- Async mode: opt-in for high throughput
- Multiple sinks: single sink by default

### 5. Graceful Degradation
- Invalid fields logged as warnings, not errors
- Sink write failures don't crash logger
- Async mode flushes queue on close

---

## Performance Characteristics

### Synchronous Mode
- **Latency**: Blocking I/O on every log call
- **Throughput**: ~68K logs/sec
- **Use case**: Low-volume logging, simple applications

### Asynchronous Mode
- **Latency**: Non-blocking (enqueue only)
- **Throughput**: ~100K logs/sec
- **Use case**: High-volume logging, performance-critical paths
- **Trade-off**: Logs may be lost if process crashes before flush

### Memory Usage
- Sync: Minimal (no buffering)
- Async: BufferSize × average entry size
- Typical: 100-500 buffer size = ~50-250KB

---

## Extension Points

### Custom Formatter
```go
type MyFormatter struct{}

func (f *MyFormatter) Format(entry *logger.Entry) ([]byte, error) {
    // Custom formatting logic
}
```

### Custom Sink
```go
type MySink struct{}

func (s *MySink) Write(data []byte) error {
    // Custom output logic
}

func (s *MySink) Close() error {
    // Cleanup
}
```

### Integration Example
```go
config := logger.Config{
    Formatter: &MyFormatter{},
    Sinks: []logger.Sink{&MySink{}},
}
log := logger.NewWithConfig(config)
```

---

## Thread Safety

- **Logger**: Safe for concurrent use (no shared mutable state)
- **Async Worker**: Channel-based synchronization
- **Sinks**: Responsibility of sink implementation
  - ConsoleSink: Safe (stdout is synchronized)
  - FileSink: Safe (os.File operations are atomic)
  - Custom: Must implement own synchronization

---

## Error Handling

### Non-Fatal Errors
- Invalid field keys → Warning to stderr, skip pair
- Sink write failure → Warning to stderr, continue to next sink
- Format failure → Warning to stderr, skip log

### Fatal Errors
- `log.Fatal()` → Logs message then calls `os.Exit(1)`
- File sink creation failure → Returns error to caller

---

## Implemented Features (v1.0)

1. **Contextual Logging**: Child loggers with inherited fields
   - `With()` method creates immutable child loggers
   - Fields merge at log time
   - Nested contexts supported

2. **Log Rotation**: Size-based automatic file rotation
   - Configurable max size and backup count
   - Thread-safe rotation
   - Automatic backup shifting

3. **Performance Benchmarks**: Comprehensive benchmark suite
   - Measures all major operations
   - Identifies performance characteristics
   - Guides optimization decisions

## Future Enhancements

1. **Sampling**: Log only N% of entries at high volume
2. **Batching**: Group multiple entries for efficient I/O
3. **Compression**: Compress old log files
4. **Remote Sinks**: HTTP, Kafka, CloudWatch, etc.
5. **Context Integration**: Extract fields from context.Context
6. **Structured Errors**: Better error type handling
