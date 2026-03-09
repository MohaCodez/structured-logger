# Structured Logger (Go)

A modular, plug-and-play structured logging library in Go featuring leveled logging, JSON output, caller tracing, and extensible output sinks.

This project is designed as a reusable infrastructure component that can be integrated into any Go application including APIs, background workers, CLI tools, and distributed services.

---

# Motivation

Logging is a fundamental part of observability in modern systems. Applications generate large volumes of runtime events which need to be structured, searchable, and consistent across services.

Traditional string-based logs are difficult to parse and analyze. Structured logging solves this by emitting logs as machine-readable events.

This logger aims to provide:

- structured JSON logs
- flexible output sinks
- clean and minimal API
- easy integration into any Go project

The goal of this project is to build a logging library that behaves like a reusable infrastructure primitive.

---

# Features

- Structured logging with JSON output
- Log levels (DEBUG, INFO, WARN, ERROR, FATAL)
- Caller tracing (file and line number)
- Context-based metadata fields
- Pluggable output sinks
- Optional asynchronous logging
- Minimal dependencies
- Designed for extensibility

---

# Architecture

The logger follows a pipeline architecture where log events pass through multiple stages before reaching their final destination.

```
Application
     │
     ▼
 Logger API
     │
     ▼
Entry Builder
     │
     ▼
Level Filter
     │
     ▼
Formatter
     │
     ▼
Sink Dispatcher
     │
     ├── Console Sink
     ├── File Sink
     └── Remote Sink
```

Each component is isolated and replaceable, allowing the system to evolve without breaking the public API.

---

# Log Entry Structure

Each log message is represented as a structured event.

Example output:

```json
{
  "timestamp": "2026-03-09T12:10:10Z",
  "level": "INFO",
  "message": "user_login",
  "user_id": 123,
  "request_id": "abc123",
  "caller": "auth/service.go:42"
}
```

Structured logs make it easy to index and query logs in centralized logging systems.

---

# Installation

Clone the repository:

```
git clone https://github.com/MohaCodez/structured-logger.git
```

Or add it as a dependency in your Go module.

---

# Quick Start

Basic usage example:

```go
logger := logger.NewDefault()

logger.Info("user_login",
    "user_id", 123,
    "ip", "10.1.2.4",
)
```

Example output:

```
{
  "timestamp":"2026-03-09T12:00:00Z",
  "level":"INFO",
  "message":"user_login",
  "user_id":123,
  "ip":"10.1.2.4"
}
```

---

# Configuration

The logger can be configured to control formatting, output sinks, and runtime behavior.

Example configuration:

```
level: INFO
formatter: json
sinks:
  - console
  - file
caller: true
async: false
```

---

# Supported Sinks

Current sinks include:

- Console output
- File logging
- Multi-sink fan-out

Future sinks may include:

- HTTP log streaming
- Message queue logging
- Cloud logging services

---

# Asynchronous Logging

For high-throughput systems, logs can optionally be processed asynchronously.

```
Application
     │
     ▼
Log Queue
     │
     ▼
Background Worker
     │
     ▼
Sink Dispatcher
```

This prevents logging from blocking application execution.

---

# Project Structure

```
logger/
formatter/
sink/
async/
middleware/
internal/
tests/
docs/
```

Each module represents a distinct part of the logging pipeline.

---

# Example Use Cases

This logger can be integrated into:

- REST APIs
- Microservices
- CLI tools
- Background workers
- Distributed systems

Because the library is modular, it can be extended for different environments and logging pipelines.

---

# Roadmap

Planned improvements include:

- asynchronous logging pipeline
- log sampling
- log rotation
- distributed tracing integration
- additional output sinks
- performance optimizations

---

# License

MIT License
