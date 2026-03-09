package main

import "github.com/MohaCodez/structured-logger/logger"

func main() {
	log := logger.New(logger.DEBUG)

	log.Debug("this is a debug message")
	log.Info("service started")
	log.Warn("high memory usage detected")
	log.Error("failed to connect to database")
	
	// Uncomment to test FATAL (will exit the program)
	// log.Fatal("critical system failure")
}
