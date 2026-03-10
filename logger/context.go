package logger

import "context"

type contextKey struct{}

var loggerKey = contextKey{}

// WithContext stores a logger in the context
func WithContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext retrieves a logger from the context.
// Returns a default logger if none is set.
func FromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(loggerKey).(*Logger); ok {
		return l
	}
	// Return default logger if none found
	return New(INFO)
}
