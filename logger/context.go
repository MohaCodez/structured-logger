package logger

import "context"

// contextKey is a private type used as a key for storing the logger in context.Context.
// Using a custom type (instead of a string) prevents collisions with other packages
// that might use context values.
type contextKey struct{}

// loggerKey is the singleton instance used to store/retrieve loggers from context.
var loggerKey = contextKey{}

// WithContext stores a logger in the given context and returns a new context.
// This is useful for passing loggers through a call chain without explicit parameters.
//
// Common pattern: Store a request-scoped logger in the request context.
//
// Example:
//
//	func handleRequest(w http.ResponseWriter, r *http.Request) {
//	    // Create a request-scoped logger with request ID
//	    requestLog := baseLogger.With("request_id", generateID())
//
//	    // Store it in the request context
//	    ctx := logger.WithContext(r.Context(), requestLog)
//
//	    // Pass context to other functions
//	    processRequest(ctx)
//	}
//
//	func processRequest(ctx context.Context) {
//	    // Retrieve the logger from context
//	    log := logger.FromContext(ctx)
//	    log.Info("processing") // Includes request_id automatically
//	}
func WithContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext retrieves a logger from the given context.
// If no logger is stored in the context, returns a default logger with INFO level.
//
// This function never returns nil, so it's always safe to use:
//
//	log := logger.FromContext(ctx)
//	log.Info("message") // Always works, even if no logger was stored
//
// Example:
//
//	func processRequest(ctx context.Context) {
//	    log := logger.FromContext(ctx)
//	    log.Info("processing request")
//
//	    // Pass context to nested functions
//	    authenticateUser(ctx)
//	}
//
//	func authenticateUser(ctx context.Context) {
//	    log := logger.FromContext(ctx)
//	    log.Info("authenticating user")
//	    // If the context has a logger with request_id, it's included here too
//	}
func FromContext(ctx context.Context) *Logger {
	// Try to retrieve the logger from context
	if l, ok := ctx.Value(loggerKey).(*Logger); ok {
		return l
	}
	// Return a default logger if none was stored
	// This ensures FromContext never returns nil
	return New(INFO)
}
