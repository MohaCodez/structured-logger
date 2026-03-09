# Performance Benchmarks

Benchmark results for the structured logger on Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz.

## Benchmark Results

| Operation | Time/op | Memory/op | Allocs/op |
|-----------|---------|-----------|-----------|
| Sync Logging | 2,078 ns | 912 B | 17 |
| Async Logging | 2,046 ns | 928 B | 18 |
| Structured Fields (3 fields) | 4,607 ns | 1,745 B | 25 |
| JSON Formatting | 2,262 ns | 880 B | 17 |
| Contextual Logging | 3,185 ns | 1,392 B | 22 |
| With Caller Tracing | 3,203 ns | 1,312 B | 24 |
| Level Filtering (filtered) | 1.7 ns | 0 B | 0 |
| Nested Contextual Logging | 3,810 ns | 1,376 B | 22 |

## Analysis

### Throughput

Based on the benchmarks:
- **Sync logging**: ~481,000 logs/sec
- **Async logging**: ~489,000 logs/sec
- **With structured fields**: ~217,000 logs/sec
- **Level filtering**: ~577 million ops/sec (essentially free)

### Key Insights

1. **Level Filtering is Extremely Fast**
   - Filtered logs cost only 1.7 ns with zero allocations
   - Always set appropriate log levels in production

2. **Async vs Sync Performance**
   - Async logging is slightly faster (2,046 ns vs 2,078 ns)
   - Similar memory footprint
   - Async provides better throughput under load

3. **Structured Fields Overhead**
   - Adding 3 fields increases time by ~2.5x
   - Memory usage increases from 912 B to 1,745 B
   - Still achieves 217K logs/sec

4. **Contextual Logging Cost**
   - Child loggers add ~1,100 ns overhead (53% increase)
   - Memory increases by ~480 B
   - Trade-off: convenience vs performance

5. **Caller Tracing Overhead**
   - Adds ~1,125 ns (54% increase)
   - Minimal memory impact
   - Disable in production for better performance

6. **Nested Context**
   - Nested child loggers add ~1,732 ns overhead
   - Reasonable cost for the convenience

## Recommendations

### Production Settings

For high-performance production systems:

```go
config := logger.Config{
    Level:        logger.INFO,        // Filter debug logs
    EnableCaller: false,              // Disable caller tracing
    Async:        true,               // Enable async mode
    BufferSize:   1000,               // Large buffer
}
```

Expected performance: **~489,000 logs/sec** with minimal overhead.

### Development Settings

For development with full debugging:

```go
config := logger.Config{
    Level:        logger.DEBUG,       // Show all logs
    EnableCaller: true,               // Enable caller tracing
    Async:        false,              // Sync for immediate output
}
```

Expected performance: **~312,000 logs/sec** (with caller tracing).

### Contextual Logging

Use contextual logging judiciously:
- Create child loggers at request/session boundaries
- Avoid creating child loggers in tight loops
- Overhead: ~1,100 ns per log with context

## Running Benchmarks

```bash
# Run all benchmarks
go test ./benchmarks -bench=. -benchmem

# Run specific benchmark
go test ./benchmarks -bench=BenchmarkSyncLogging -benchmem

# Run with more iterations
go test ./benchmarks -bench=. -benchtime=10s
```

## Comparison with Standard Library

The structured logger provides:
- **Structured fields**: Native support (vs manual formatting)
- **Multiple sinks**: Built-in fan-out
- **Async mode**: Non-blocking logging
- **Contextual logging**: Inherited fields

Trade-off: ~2,000 ns/op vs ~500 ns/op for basic `log.Printf()`, but with significantly more features.

## Memory Profile

Allocations per operation:
- **Minimal**: 17 allocs (basic logging)
- **With fields**: 25 allocs (3 fields)
- **With context**: 22 allocs (inherited fields)

Most allocations are for:
1. Entry struct creation
2. Field map allocation
3. JSON marshaling
4. Timestamp formatting

## Optimization Tips

1. **Use appropriate log levels** - Level filtering is free
2. **Disable caller in production** - Saves ~1,100 ns/op
3. **Enable async mode** - Better throughput under load
4. **Reuse child loggers** - Don't create in loops
5. **Limit field count** - Each field adds overhead

---

*Benchmarks run on: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz*
