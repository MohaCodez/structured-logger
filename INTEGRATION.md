# Integration Guide

## Getting Started

This guide shows how to integrate the structured logger into your Go application.

---

## Installation

```bash
go get github.com/MohaCodez/structured-logger
```

---

## Basic Integration

### 1. Simple Application

```go
package main

import (
    "github.com/MohaCodez/structured-logger/logger"
)

func main() {
    // Create logger
    log := logger.New(logger.INFO)
    defer log.Close()

    // Use throughout application
    log.Info("application_started")
    
    // Your application code
    processRequests(log)
}

func processRequests(log *logger.Logger) {
    log.Debug("processing_started")
    // ... business logic
    log.Info("request_completed", "duration_ms", 45)
}
```

---

## Production Configuration

### 2. Production Setup

```go
package main

import (
    "os"
    
    "github.com/MohaCodez/structured-logger/formatter"
    "github.com/MohaCodez/structured-logger/logger"
    "github.com/MohaCodez/structured-logger/sink"
)

func main() {
    log := setupLogger()
    defer log.Close()

    log.Info("server_starting", "port", 8080)
    // ... rest of application
}

func setupLogger() *logger.Logger {
    // Determine environment
    env := os.Getenv("ENV")
    
    // Configure based on environment
    config := logger.DefaultConfig()
    
    if env == "production" {
        config.Level = logger.INFO
        config.EnableCaller = false
        config.Async = true
        config.BufferSize = 500
        
        // Add file sink
        fileSink, err := sink.NewFileSink("/var/log/app/app.log")
        if err != nil {
            panic(err)
        }
        config.Sinks = []logger.Sink{
            sink.NewConsoleSink(),
            fileSink,
        }
    } else {
        // Development settings
        config.Level = logger.DEBUG
        config.EnableCaller = true
        config.Async = false
    }
    
    config.Formatter = formatter.NewJSONFormatter()
    
    return logger.NewWithConfig(config)
}
```

---

## HTTP Server Integration

### 3. REST API with Middleware

```go
package main

import (
    "net/http"
    "time"
    
    "github.com/MohaCodez/structured-logger/logger"
)

type Server struct {
    log *logger.Logger
}

func main() {
    log := logger.New(logger.INFO)
    defer log.Close()
    
    server := &Server{log: log}
    
    http.HandleFunc("/api/users", server.loggingMiddleware(server.handleUsers))
    
    log.Info("server_listening", "port", 8080)
    http.ListenAndServe(":8080", nil)
}

func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        s.log.Info("request_started",
            "method", r.Method,
            "path", r.URL.Path,
            "remote_addr", r.RemoteAddr,
        )
        
        next(w, r)
        
        s.log.Info("request_completed",
            "method", r.Method,
            "path", r.URL.Path,
            "duration_ms", time.Since(start).Milliseconds(),
        )
    }
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
    s.log.Debug("fetching_users")
    
    // Business logic
    users := fetchUsers()
    
    if len(users) == 0 {
        s.log.Warn("no_users_found")
    }
    
    s.log.Info("users_fetched", "count", len(users))
    
    // Send response
    w.Write([]byte("users"))
}

func fetchUsers() []string {
    return []string{"alice", "bob"}
}
```

---

## Dependency Injection

### 4. Passing Logger Through Application

```go
package main

import (
    "github.com/MohaCodez/structured-logger/logger"
)

type App struct {
    log     *logger.Logger
    service *UserService
}

type UserService struct {
    log *logger.Logger
}

func main() {
    log := logger.New(logger.INFO)
    defer log.Close()
    
    app := NewApp(log)
    app.Run()
}

func NewApp(log *logger.Logger) *App {
    return &App{
        log:     log,
        service: NewUserService(log),
    }
}

func (a *App) Run() {
    a.log.Info("app_started")
    a.service.CreateUser("alice")
}

func NewUserService(log *logger.Logger) *UserService {
    return &UserService{log: log}
}

func (s *UserService) CreateUser(name string) {
    s.log.Info("creating_user", "name", name)
    
    // Business logic
    if err := validateUser(name); err != nil {
        s.log.Error("user_validation_failed",
            "name", name,
            "error", err.Error(),
        )
        return
    }
    
    s.log.Info("user_created", "name", name)
}

func validateUser(name string) error {
    return nil
}
```

---

## Error Handling

### 5. Logging Errors with Context

```go
package main

import (
    "errors"
    
    "github.com/MohaCodez/structured-logger/logger"
)

func main() {
    log := logger.New(logger.INFO)
    defer log.Close()
    
    if err := processPayment(log, "txn_123", 99.99); err != nil {
        log.Error("payment_failed",
            "transaction_id", "txn_123",
            "error", err.Error(),
        )
    }
}

func processPayment(log *logger.Logger, txnID string, amount float64) error {
    log.Info("payment_processing",
        "transaction_id", txnID,
        "amount", amount,
    )
    
    // Simulate error
    if amount > 100 {
        log.Warn("high_value_transaction",
            "transaction_id", txnID,
            "amount", amount,
        )
    }
    
    // Business logic
    if err := chargeCard(amount); err != nil {
        log.Error("card_charge_failed",
            "transaction_id", txnID,
            "amount", amount,
            "error", err.Error(),
        )
        return err
    }
    
    log.Info("payment_completed",
        "transaction_id", txnID,
        "amount", amount,
    )
    
    return nil
}

func chargeCard(amount float64) error {
    return errors.New("insufficient funds")
}
```

---

## Background Workers

### 6. Worker Pool with Logging

```go
package main

import (
    "sync"
    "time"
    
    "github.com/MohaCodez/structured-logger/logger"
)

func main() {
    log := logger.New(logger.INFO)
    defer log.Close()
    
    jobs := make(chan int, 100)
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go worker(log, i, jobs, &wg)
    }
    
    // Send jobs
    for j := 0; j < 20; j++ {
        jobs <- j
    }
    close(jobs)
    
    wg.Wait()
    log.Info("all_workers_completed")
}

func worker(log *logger.Logger, id int, jobs <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    log.Info("worker_started", "worker_id", id)
    
    for job := range jobs {
        log.Debug("processing_job",
            "worker_id", id,
            "job_id", job,
        )
        
        time.Sleep(100 * time.Millisecond)
        
        log.Info("job_completed",
            "worker_id", id,
            "job_id", job,
        )
    }
    
    log.Info("worker_stopped", "worker_id", id)
}
```

---

## Testing with Logger

### 7. Mock Logger for Tests

```go
package main

import (
    "testing"
    
    "github.com/MohaCodez/structured-logger/logger"
)

// Mock sink for testing
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

func TestUserService(t *testing.T) {
    // Create logger with test sink
    sink := &testSink{}
    log := logger.NewWithSinks(
        logger.DEBUG,
        &logger.DefaultFormatter{},
        []logger.Sink{sink},
    )
    
    service := NewUserService(log)
    service.CreateUser("alice")
    
    // Verify logs
    if len(sink.logs) == 0 {
        t.Error("expected logs to be written")
    }
}
```

---

## Best Practices

### 1. Logger Lifecycle
- Create logger at application startup
- Pass logger to components via dependency injection
- Always defer `log.Close()` in main

### 2. Log Levels
- **DEBUG**: Development details, verbose
- **INFO**: Normal operations, user actions
- **WARN**: Unexpected but handled conditions
- **ERROR**: Errors that need attention
- **FATAL**: Critical errors, application exits

### 3. Structured Fields
- Use consistent field names across application
- Include request IDs for tracing
- Add user context when available
- Include timing information

### 4. Performance
- Use async mode for high-throughput applications
- Disable caller tracing in production
- Set appropriate log level (INFO or WARN in prod)

### 5. Error Context
- Always include error message in fields
- Add relevant context (IDs, values)
- Log at appropriate level

---

## Common Patterns

### Request ID Tracking
```go
func handleRequest(log *logger.Logger, requestID string) {
    log.Info("request_started", "request_id", requestID)
    // ... processing
    log.Info("request_completed", "request_id", requestID)
}
```

### Timing Operations
```go
start := time.Now()
// ... operation
log.Info("operation_completed", "duration_ms", time.Since(start).Milliseconds())
```

### Conditional Logging
```go
if result.Count > threshold {
    log.Warn("threshold_exceeded", "count", result.Count, "threshold", threshold)
}
```

---

## Troubleshooting

### Logs Not Appearing
- Check log level configuration
- Verify sinks are properly initialized
- Ensure `Close()` is called (flushes async queue)

### Performance Issues
- Enable async mode
- Reduce log volume (increase level threshold)
- Disable caller tracing

### File Sink Issues
- Check file permissions
- Verify directory exists
- Ensure disk space available
