package sink

// Sink is an interface for log output destinations.
// Implement this interface to send logs to custom destinations like:
//   - Databases (PostgreSQL, MongoDB, etc.)
//   - Message queues (Kafka, RabbitMQ, etc.)
//   - Cloud services (AWS CloudWatch, GCP Logging, etc.)
//   - Network endpoints (HTTP, TCP, UDP, etc.)
//
// The library provides ConsoleSink, FileSink, and RotatingFileSink by default.
//
// Example custom sink:
//
//	type HTTPSink struct {
//	    url string
//	}
//
//	func (s *HTTPSink) Write(data []byte) error {
//	    _, err := http.Post(s.url, "application/json", bytes.NewReader(data))
//	    return err
//	}
//
//	func (s *HTTPSink) Close() error {
//	    return nil // No resources to clean up
//	}
type Sink interface {
	// Write outputs the formatted log data to the destination.
	// The data is already formatted (e.g., as JSON) by the formatter.
	// Returns an error if the write operation fails.
	Write(data []byte) error

	// Close releases any resources held by the sink.
	// This is called when the logger is closed.
	// Examples: closing file handles, network connections, database connections.
	Close() error
}
