package main

import (
	"fmt"
	"time"

	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

func main() {
	jsonFormatter := formatter.NewJSONFormatter()
	consoleSink := sink.NewConsoleSink()

	fmt.Println("=== Synchronous Logging ===")
	syncStart := time.Now()
	syncLog := logger.NewWithCaller(
		logger.INFO,
		jsonFormatter,
		[]logger.Sink{consoleSink},
		false,
	)
	for i := 0; i < 1000; i++ {
		syncLog.Info("test", "iteration", i)
	}
	syncLog.Close()
	syncDuration := time.Since(syncStart)

	fmt.Printf("\n=== Asynchronous Logging ===\n")
	asyncStart := time.Now()
	asyncLog := logger.NewAsync(
		logger.INFO,
		jsonFormatter,
		[]logger.Sink{consoleSink},
		false,
		500, // buffer size
	)
	for i := 0; i < 1000; i++ {
		asyncLog.Info("test", "iteration", i)
	}
	asyncLog.Close()
	asyncDuration := time.Since(asyncStart)

	fmt.Printf("\n=== Performance Comparison ===\n")
	fmt.Printf("Sync:  %v\n", syncDuration)
	fmt.Printf("Async: %v\n", asyncDuration)
	fmt.Printf("Speedup: %.2fx\n", float64(syncDuration)/float64(asyncDuration))
}
