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

	// Join the provided arguments into a single message string.
	message := joinMessages(" ---> ", args...)

	// Combine the new message with the original error's message.
	fullMessage := joinMessages(" ---> ", message, err.Error())

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
// Logs the error using the appropriate logger based on the current logging level.
func Log(err error, req interface{}, msgs ...interface{}) {
	if err == nil {
		return
	}

	// Asynchronously handle the error logging to prevent blocking.
	go logError(err, req, msgs...)
}

// logError is a helper function that logs the error using the appropriate logger.
// It constructs a combined message and uses the configured loggers based on the current logging level.
func logError(err error, req interface{}, msgs ...interface{}) {
	// Join all provided messages to create a unified error message.
	message := joinMessages(" ---> ", msgs...)
	errorPath := slog.String("Error Path", err.Error())

	// Retrieve the logger based on the current logging level.
	getLogger(message, errorPath, slog.Any("request", req))
}

// getLogger returns the appropriate logger based on the current logging level.
func getLogger(msg string, args ...any) {
	switch currentLevel {
	case LevelLocal:
		slogTextLogger.Error(msg, args...)
	case LevelStaging, LevelMaster:
		if fileLogger != nil {
			fileLogger.Error(msg, args...)
		}
		slogJSONLogger.Error(msg, args...)
	default:
		slogJSONLogger.Error(msg, args...)
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
