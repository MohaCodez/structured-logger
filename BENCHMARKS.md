# Performance Benchmarks

A comprehensive performance analysis of the structured logging library.

---

## Benchmark Environment

All benchmarks were executed on the following system configuration:

| Component | Specification |
|-----------|---------------|
| **CPU** | Intel Core i7 (9th Generation) |
| **RAM** | 8 GB |
| **Operating System** | Ubuntu 22.04 LTS |
| **Architecture** | amd64 |
| **Go Version** | go1.25.7 linux/amd64 |

**Note:** Benchmark results may vary depending on:
- CPU performance and clock speed
- Available memory and system load
- Operating system scheduling policies
- Go runtime version and compiler optimizations

For consistent results, run benchmarks on an idle system with minimal background processes.

---

## Running Benchmarks

Developers can reproduce these benchmarks locally using the following commands:

### Run All Benchmarks
```bash
go test -bench=. -benchmem ./...
```

### Run Specific Benchmark
```bash
go test -bench=BenchmarkSyncLogging -benchmem ./benchmarks
```

### Run with Extended Duration
```bash
go test -bench=. -benchmem -benchtime=10s ./benchmarks
```

### Command Options
- **`-bench=.`** - Runs all benchmark functions matching the pattern
- **`-benchmem`** - Reports memory allocation statistics
- **`-benchtime=Ns`** - Runs each benchmark for N seconds (default: 1s)

---

## Understanding Benchmark Metrics

Go benchmarks report three key metrics for each operation:

### ns/op (Nanoseconds per Operation)
Time required to complete one logging operation, measured in nanoseconds.

**Lower is better.** Indicates the latency of a single log call.

### B/op (Bytes per Operation)
Number of bytes allocated in memory per operation.

**Lower is better.** High memory allocation increases garbage collection pressure.

### allocs/op (Allocations per Operation)
Number of distinct memory allocations per operation.

**Lower is better.** Each allocation requires GC tracking. In high-throughput systems, excessive allocations can degrade performance significantly.

**Why This Matters for Logging:**

Logging libraries are often called thousands of times per second. Even small per-operation costs multiply quickly:
- A 2,000 ns/op logger handling 10,000 logs/sec consumes ~20ms of CPU time
- High allocation rates trigger frequent garbage collection pauses
- Memory-efficient logging reduces GC overhead in production systems

---

## Estimated Throughput

Benchmark latency (ns/op) can be converted to approximate throughput using:

```
logs_per_second ≈ 1,000,000,000 / ns_per_op
```

### Example Calculation

If a benchmark reports:
```
BenchmarkSyncLogging    1200 ns/op
```

Then approximate throughput is:
```
1,000,000,000 / 1200 ≈ 833,000 logs/second
```

This provides an intuitive understanding of logging capacity under sustained load.

**Note:** Actual throughput depends on:
- Log message complexity
- Number of structured fields
- Sink I/O performance
- System contention

---

## Benchmark Scenarios

### BenchmarkSyncLogging
Measures the cost of a synchronous logging call with minimal fields.

**Pipeline:**
```
Application → Logger → Formatter → Sink
```

**What it tests:** End-to-end latency of a blocking log operation.

---

### BenchmarkAsyncLogging
Measures the cost of enqueueing logs when async logging is enabled.

**Pipeline:**
```
Application → Queue → [returns immediately]
                ↓
         Worker → Formatter → Sink
```

**What it tests:** Non-blocking enqueue latency. Actual formatting and I/O happen asynchronously.

---

### BenchmarkStructuredFields
Measures overhead introduced by structured metadata fields (3 fields).

**Pipeline:**
```
Application → Logger → Field Parsing → Formatter → Sink
```

**What it tests:** Cost of adding contextual key-value pairs to log entries.

---

### BenchmarkJSONFormatting
Measures JSON serialization performance.

**What it tests:** Formatter overhead converting structured entries to JSON bytes.

---

### BenchmarkContextualLogging
Measures performance of child loggers with inherited fields.

**What it tests:** Cost of context propagation and field merging.

---

### BenchmarkWithCaller
Measures overhead of capturing caller information (file:line).

**What it tests:** Runtime stack inspection cost.

---

### BenchmarkLevelFiltering
Measures cost of filtering logs below the configured level.

**What it tests:** Early-exit performance when logs are discarded.

---

### BenchmarkNestedContextualLogging
Measures performance of deeply nested child loggers.

**What it tests:** Cost of multiple levels of context inheritance.

---

## Benchmark Results

Results from the configured benchmark environment:

| Operation | Time/op | Memory/op | Allocs/op | Throughput |
|-----------|---------|-----------|-----------|------------|
| Sync Logging | 2,078 ns | 912 B | 17 | ~481K logs/sec |
| Async Logging | 2,046 ns | 928 B | 18 | ~489K logs/sec |
| Structured Fields (3 fields) | 4,607 ns | 1,745 B | 25 | ~217K logs/sec |
| JSON Formatting | 2,262 ns | 880 B | 17 | ~442K logs/sec |
| Contextual Logging | 3,185 ns | 1,392 B | 22 | ~314K logs/sec |
| With Caller Tracing | 3,203 ns | 1,312 B | 24 | ~312K logs/sec |
| Level Filtering (filtered) | 1.7 ns | 0 B | 0 | ~577M ops/sec |
| Nested Contextual Logging | 3,810 ns | 1,376 B | 22 | ~262K logs/sec |

---

## Analysis

### Key Insights

**1. Level Filtering is Extremely Fast**
- Filtered logs cost only 1.7 ns with zero allocations
- Always set appropriate log levels in production to skip unnecessary work

**2. Async vs Sync Performance**
- Async logging provides ~1.5% better throughput (489K vs 481K logs/sec)
- Similar memory footprint (928 B vs 912 B)
- Async mode shines under high concurrency and I/O-bound sinks

**3. Structured Fields Overhead**
- Adding 3 fields increases latency by ~2.5x (4,607 ns vs 2,078 ns)
- Memory usage nearly doubles (1,745 B vs 912 B)
- Trade-off: rich context vs raw speed

**4. Contextual Logging Cost**
- Child loggers add ~1,100 ns overhead (53% increase)
- Memory increases by ~480 B
- Reasonable cost for avoiding repeated field passing

**5. Caller Tracing Overhead**
- Adds ~1,125 ns (54% increase) due to runtime stack inspection
- Minimal memory impact
- Recommended: disable in production, enable in development

**6. Nested Context Performance**
- Nested child loggers add ~1,732 ns overhead (83% increase)
- Acceptable for request-scoped logging patterns

---

## Recommendations

### Production Configuration

For high-performance production systems:

```go
config := logger.Config{
    Level:        logger.INFO,        // Filter debug logs
    EnableCaller: false,              // Disable caller tracing
    Async:        true,               // Enable async mode
    BufferSize:   1000,               // Large buffer for bursts
}
```

**Expected performance:** ~489,000 logs/sec with minimal overhead.

---

### Development Configuration

For development with full debugging:

```go
config := logger.Config{
    Level:        logger.DEBUG,       // Show all logs
    EnableCaller: true,               // Enable caller tracing
    Async:        false,              // Sync for immediate output
}
```

**Expected performance:** ~312,000 logs/sec (with caller tracing).

---

### Contextual Logging Guidelines

Use contextual logging judiciously:
- Create child loggers at request/session boundaries
- Avoid creating child loggers in tight loops
- Overhead: ~1,100 ns per log with inherited context

---

## Comparison Guidance

Developers may compare these results with other popular Go logging libraries:

- **[zap](https://github.com/uber-go/zap)** - High-performance structured logger
- **[zerolog](https://github.com/rs/zerolog)** - Zero-allocation JSON logger
- **[slog](https://pkg.go.dev/log/slog)** - Standard library structured logger (Go 1.21+)

### Performance vs Design Trade-offs

This library prioritizes:
- **Modular architecture** - Pluggable formatters and sinks
- **Extensibility** - Easy to add custom components
- **Clarity** - Readable code over micro-optimizations

While raw speed is important, this project balances performance with:
- Clean separation of concerns
- Testability and maintainability
- Flexibility for diverse use cases

**Benchmark Context:**

Some libraries achieve sub-1000 ns/op by:
- Pre-allocating buffers
- Using sync.Pool for object reuse
- Specialized zero-allocation JSON encoders

This library achieves competitive performance (~2,000 ns/op) while maintaining a simple, extensible design suitable for learning and customization.

---

## Memory Profile

Allocations per operation breakdown:

| Component | Allocations |
|-----------|-------------|
| **Minimal logging** | 17 allocs |
| **With 3 fields** | 25 allocs |
| **With context** | 22 allocs |

Most allocations occur in:
1. Entry struct creation
2. Field map allocation
3. JSON marshaling
4. Timestamp formatting

---

## Optimization Tips

1. **Use appropriate log levels** - Level filtering is essentially free (1.7 ns)
2. **Disable caller in production** - Saves ~1,100 ns/op
3. **Enable async mode** - Better throughput under load
4. **Reuse child loggers** - Don't create in hot paths
5. **Limit field count** - Each field adds parsing and serialization overhead

---

## Running Custom Benchmarks

To benchmark specific scenarios:

```bash
# Benchmark with CPU profiling
go test -bench=. -benchmem -cpuprofile=cpu.prof ./benchmarks

# Benchmark with memory profiling
go test -bench=. -benchmem -memprofile=mem.prof ./benchmarks

# Analyze profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

---

*Benchmarks executed on: Intel Core i7 (9th Gen), 8GB RAM, Ubuntu 22.04 LTS, Go 1.25.7*
