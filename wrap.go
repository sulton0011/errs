package errs

import (
	"fmt"
	"log/slog"
	"strings"
)

// Wrap adds context to an existing error by wrapping it with additional messages.
// It accepts an error and variadic arguments to append to the error message chain.
// Returns a new error with the combined messages and original error context.
func Wrap(err error, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return wrap(err, args...)
}

func wrap(err error, args ...interface{}) error {
	// Join the provided arguments into a single message string.
	message := joinMessages(" ---> ", args...)

	// Combine the new message with the original error's message.
	fullMessage := joinMessages(" ---> ", message, Unwrap(err))

	return &errorString{
		message: fullMessage,
		origErr: err.Error(),
	}
}

// Wrapf formats a string message using the provided format and arguments,
// then wraps it around an existing error to add context.
// Returns a new error with the formatted message and original error context.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	// Use fmt.Sprintf to format the message, then wrap the error.
	return Wrap(err, fmt.Sprintf(format, args...))
}

// Log asynchronously logs an error with additional context messages and a request object.
// If the error is nil, it does nothing.
// Logs the error using the appropriate logger based on the current logging level and type.
func Log(err error, req interface{}, msgs ...interface{}) {
	if err == nil {
		return
	}

	// Asynchronously handle the error logging to prevent blocking.
	go logError(err, req, msgs...)
}

// Unwrap retrieves the original message from a wrapped error.
func Unwrap(err error) string {
	u, ok := err.(*errorString)
	if !ok {
		return err.Error()
	}

	return u.message
}

// logError is a helper function that logs the error using the appropriate logger.
// It constructs a combined message and uses the configured loggers based on the current logging level and type.
func logError(err error, req interface{}, msgs ...interface{}) {
	// Join all provided messages to create a unified error message.
	message := joinMessages(" ---> ", msgs...)
	errorPath := slog.String("Error Path", Unwrap(err))

	// Retrieve the logger based on the current logging level and type.
	getLogger(message, errorPath, slog.Any("request", req))
}

// getLogger uses the appropriate logger based on the current logging level and type.
func getLogger(msg string, args ...any) {
	for _, slog := range slogLoggers {
		go slog.Error(msg, args...)
	}
}

// joinMessages concatenates multiple messages using the specified separator.
// It skips nil values and ensures a clean, unified message string.
func joinMessages(separator string, messages ...interface{}) string {
	var builder strings.Builder
	for _, message := range messages {
		if message != nil {
			if builder.Len() > 0 {
				builder.WriteString(separator) // Append separator between messages.
			}
			builder.WriteString(fmt.Sprint(message)) // Convert and append message to builder.
		}
	}
	return builder.String()
}
