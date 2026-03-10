package main

import (
	"fmt"

	"github.com/MohaCodez/structured-logger/formatter"
	"github.com/MohaCodez/structured-logger/logger"
	"github.com/MohaCodez/structured-logger/sink"
)

func main() {
	jsonFormatter := formatter.NewJSONFormatter()
	consoleSink := sink.NewConsoleSink()

	fmt.Println("=== Without Caller ===")
	logWithoutCaller := logger.NewWithCaller(
		logger.INFO,
		jsonFormatter,
		[]logger.Sink{consoleSink},
		false, // caller disabled
	)
	logWithoutCaller.Info("test_message", "key", "value")
	logWithoutCaller.Close()

	fmt.Println("\n=== With Caller ===")
	logWithCaller := logger.NewWithCaller(
		logger.INFO,
		jsonFormatter,
		[]logger.Sink{consoleSink},
		true, // caller enabled
	)
	logWithCaller.Info("test_message", "key", "value")
	logWithCaller.Close()
}
