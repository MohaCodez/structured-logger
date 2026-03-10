package logger

// Level represents the severity of a log message.
// Lower numeric values indicate more verbose logging.
type Level int

const (
	// DEBUG is the most verbose level, used for detailed debugging information.
	// Typically disabled in production environments.
	DEBUG Level = iota // iota starts at 0 and increments for each constant

	// INFO is for general informational messages about application operations.
	// This is the recommended minimum level for production.
	INFO

	// WARN indicates potentially harmful situations that should be reviewed.
	// The application continues to function normally.
	WARN

	// ERROR indicates error conditions that should be investigated.
	// The application may continue but functionality may be impaired.
	ERROR

	// FATAL indicates critical errors that require the application to exit.
	// Calling Fatal() will log the message and terminate the program.
	FATAL
)

// String converts a Level to its string representation.
// This is used when formatting log entries for output.
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}
