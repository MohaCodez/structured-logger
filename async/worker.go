package async

import (
	"fmt"
	"os"
	"sync/atomic"
)

// LogEntry represents a single log entry waiting to be written.
// It contains the formatted log data and the list of sinks to write to.
type LogEntry struct {
	Data  []byte // The formatted log data (e.g., JSON bytes)
	Sinks []Sink // The output destinations for this log entry
}

// Sink is the interface for output destinations.
// This is the same interface as logger.Sink, but defined here to avoid
// circular dependencies between packages.
type Sink interface {
	Write(data []byte) error
}

// Worker is a background goroutine that processes log entries asynchronously.
// It receives log entries through a buffered channel and writes them to sinks.
//
// Benefits of async logging:
//   - Non-blocking: Log calls return immediately without waiting for I/O
//   - Better throughput: Batches I/O operations in the background
//   - Decouples logging from application logic
//
// Trade-offs:
//   - Logs may be lost if the process crashes before the buffer is flushed
//   - Uses more memory (buffer size × average log entry size)
//   - Slightly more complex (requires proper shutdown with Stop())
type Worker struct {
	queue        chan LogEntry // Buffered channel for queuing log entries
	done         chan struct{} // Signals when the worker goroutine has finished
	dropOnFull   bool          // If true, drop logs when buffer is full; if false, block
	droppedCount uint64        // Counter for dropped logs (atomic, safe for concurrent access)
	errorHandler func(error)   // Called when a sink write fails
}

// NewWorker creates and starts a new async worker.
// The worker runs in a background goroutine and processes log entries from the queue.
//
// Parameters:
//   - bufferSize: Size of the queue (e.g., 100 means up to 100 logs can be queued)
//   - dropOnFull: If true, drop logs when buffer is full; if false, block the caller
//   - errorHandler: Function to call when a sink write fails (nil uses default stderr handler)
//
// The worker goroutine starts immediately and runs until Stop() is called.
//
// Example:
//
//	worker := async.NewWorker(500, false, nil)
//	defer worker.Stop() // Always stop the worker to flush the queue
func NewWorker(bufferSize int, dropOnFull bool, errorHandler func(error)) *Worker {
	// Use default error handler if none provided
	if errorHandler == nil {
		errorHandler = func(err error) {
			fmt.Fprintf(os.Stderr, "async worker: failed to write to sink: %v\n", err)
		}
	}

	// Create the worker
	w := &Worker{
		queue:        make(chan LogEntry, bufferSize), // Buffered channel
		done:         make(chan struct{}),             // Unbuffered channel for signaling
		dropOnFull:   dropOnFull,
		errorHandler: errorHandler,
	}

	// Start the background goroutine
	go w.run()

	return w
}

// run is the main loop of the worker goroutine.
// It continuously reads log entries from the queue and writes them to sinks.
//
// The loop exits when the queue channel is closed (by Stop()).
// After processing all remaining entries, it signals completion via the done channel.
func (w *Worker) run() {
	// Range over the channel: reads entries until the channel is closed
	for entry := range w.queue {
		// Write to all sinks for this entry
		for _, sink := range entry.Sinks {
			if err := sink.Write(entry.Data); err != nil {
				// Don't stop processing if one sink fails
				// Just report the error and continue
				w.errorHandler(err)
			}
		}
	}

	// Queue is closed and all entries processed
	// Signal that we're done
	close(w.done)
}

// Enqueue adds a log entry to the queue for background processing.
// This is called by the logger when async mode is enabled.
//
// Behavior depends on dropOnFull setting:
//
// If dropOnFull is false (default):
//   - Blocks if the queue is full (provides backpressure)
//   - Guarantees no logs are lost (unless process crashes)
//   - May slow down the caller if logs are generated faster than they can be written
//
// If dropOnFull is true:
//   - Never blocks the caller
//   - Drops the log if the queue is full
//   - Increments droppedCount (check with DroppedCount())
//   - Use when logging must never impact application performance
//
// Parameters:
//   - data: The formatted log data (already converted to bytes by the formatter)
//   - sinks: The list of output destinations for this log
func (w *Worker) Enqueue(data []byte, sinks []Sink) {
	entry := LogEntry{Data: data, Sinks: sinks}

	if w.dropOnFull {
		// Non-blocking mode: drop if buffer is full
		select {
		case w.queue <- entry:
			// Successfully queued
		default:
			// Queue is full, drop the log and increment counter
			// atomic.AddUint64 is thread-safe (can be called by multiple goroutines)
			atomic.AddUint64(&w.droppedCount, 1)
		}
	} else {
		// Blocking mode: wait for buffer space (backpressure)
		// This line blocks if the queue is full until space becomes available
		w.queue <- entry
	}
}

// DroppedCount returns the number of log entries that were dropped
// because the queue was full. Only relevant when dropOnFull is true.
//
// This uses atomic operations to safely read the counter from any goroutine.
//
// Example:
//
//	worker := async.NewWorker(10, true, nil) // Small buffer, drop on full
//	// ... lots of logging ...
//	if dropped := worker.DroppedCount(); dropped > 0 {
//	    fmt.Printf("Warning: %d logs were dropped\n", dropped)
//	}
func (w *Worker) DroppedCount() uint64 {
	// atomic.LoadUint64 safely reads the counter without a mutex
	return atomic.LoadUint64(&w.droppedCount)
}

// Stop gracefully shuts down the worker.
// It closes the queue (no more entries can be added) and waits for
// all queued entries to be processed before returning.
//
// Process:
//  1. Close the queue channel (signals the worker to stop after processing remaining entries)
//  2. Wait for the done channel (blocks until worker finishes all pending entries)
//
// Always call Stop() before the application exits to ensure logs aren't lost:
//
//	worker := async.NewWorker(100, false, nil)
//	defer worker.Stop()
//
// This is automatically called by logger.Close().
func (w *Worker) Stop() {
	// Close the queue - no more entries can be added
	// The worker goroutine will process all remaining entries and then exit
	close(w.queue)

	// Wait for the worker to finish processing all entries
	// This blocks until the worker closes the done channel
	<-w.done
}
