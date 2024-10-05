package errs

import (
	"fmt"
	"log/slog"
)

// Wrap adds context to an existing error by wrapping it with additional messages.
// It accepts an error and variadic arguments to append to the error message chain.
// Returns a new error with the combined messages and original error context.
//
// Parameters:
//   - err: The error to wrap. If this is nil, the function returns nil.
//   - args: Variadic arguments representing the additional messages to append to the error message chain.
//
// Returns:
//   - A new error with the combined messages and original error context.
//     If the provided error is nil, the function returns nil.
func Wrap(err error, args ...any) error {
	if err == nil {
		return nil
	}

	return wrap(err, args...)
}

// WrapF formats a string message using the provided format and arguments,
// then wraps it around an existing error to add context.
// If the error is nil, it returns nil.
//
// Parameters:
//   - err: The error to wrap. If this is nil, the function returns nil.
//   - format: A string format to be used with fmt.Sprintf to create the message.
//   - args: Variadic arguments representing the values to be formatted into the message.
//
// Returns:
//   - A new error with the formatted message and original error context.
//     If the provided error is nil, the function returns nil.
func WrapF(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	// Use fmt.Sprintf to format the message, then wrap the error.
	return wrap(err, fmt.Sprintf(format, args...))
}

// wrap adds context to an existing error by wrapping it with additional messages.
// It accepts an error and variadic arguments to append to the error message chain.
// Returns a new error with the combined messages and original error context.
//
// Parameters:
//   - err: The error to wrap. If this is nil, the function returns nil.
//   - args: Variadic arguments representing the additional messages to append to the error message chain.
//
// Returns:
//   - A new error with the combined messages and original error context.
//     If the provided error is nil, the function returns nil.
func wrap(err error, args ...any) error {
	// Join the provided arguments into a single message string.
	message := JoinMsg(separator, args...)

	// Combine the new message with the original error's message.
	fullMessage := JoinMsg(separator, message, Unwrap(err))

	return &errorString{
		message: fullMessage,
		origErr: err.Error(),
	}
}

// Unwrap retrieves the original message from a wrapped error.
//
// The function unwraps the error by checking if it is of type *errorString.
// If the error is wrapped, it returns the original message from the wrapped error.
// If the error is not wrapped, it returns the error message as is.
//
// Parameters:
// - err: The error to unwrap.
//
// Returns:
//   - A string representing the original message if the error is wrapped.
//     If the error is not wrapped, it returns the error message as is.
//     If the provided error is nil, it returns an empty string.
func Unwrap(err error) string {
	u, ok := err.(*errorString)
	if !ok {
		return err.Error()
	}

	return u.message
}

// UnwrapE unwraps a wrapped error and returns the original error message as a new error.
// If the provided error is not wrapped, it returns the error as is.
// If the original message is empty, it returns nil.
//
// Parameters:
//   - err: The error to unwrap. This can be nil.
//
// Returns:
//   - An error representing the original message if the error is wrapped.
//     If the error is not wrapped, it returns the error as is.
//     If the original message is empty, it returns nil.
func UnwrapE(err error) error {
	if e, ok := err.(*errorString); ok && e.message != "" {
		return New(e.message)
	}
	return err
}

// Log asynchronously logs an error with additional context messages and a request object.
// If the error is nil, it does nothing.
//
// The function logs the error using the appropriate logger based on the current logging level and type.
// It constructs a combined message by joining all provided messages with the separator separator.
// The original error message is logged as the "Error Path" field in the log entry.
// The request object is logged as the "request" field in the log entry.
//
// Parameters:
//   - err: The error to log. If this is nil, the function returns without doing anything.
//   - req: The request object associated with the error. This can be of any type.
//   - msgs: Variadic arguments representing additional messages to include in the log entry.
//
// Returns:
//   - This function does not return any value.
func Log(err error, req any, msgs ...any) {
	if err == nil {
		return
	}

	// Asynchronously handle the error logging to prevent blocking.
	go logError(err, req, msgs...)
}

// logError asynchronously logs an error with additional context messages and a request object.
// If the error is nil, it does nothing.
//
// The function logs the error using the appropriate logger based on the current logging level and type.
// It constructs a combined message by joining all provided messages with the separator separator.
// The original error message is logged as the "Error Path" field in the log entry.
// The request object is logged as the "request" field in the log entry.
//
// Parameters:
//   - err: The error to log. If this is nil, the function returns without doing anything.
//   - req: The request object associated with the error. This can be of any type.
//   - msgs: Variadic arguments representing additional messages to include in the log entry.
//
// Returns:
//   - This function does not return any value.
func logError(err error, req any, msgs ...any) {
	// Join all provided messages to create a unified error message.
	message := JoinMsg(separator, msgs...)
	errorPath := slog.String("Error Path", Unwrap(err))

	// Retrieve the logger based on the current logging level and type.
	getLogger(message, errorPath, slog.Any("request", req))
}

// getLogger retrieves and logs an error message with additional context using the appropriate logger.
// It iterates through a list of sloggers and asynchronously logs the error message with the provided arguments.
//
// Parameters:
//   - msg: A string representing the error message to be logged.
//   - args: Variadic arguments representing additional context to be logged alongside the error message.
//
// Returns:
//   - This function does not return any value. It logs the error message asynchronously using the slogLoggers.
func getLogger(msg string, args ...any) {
	for _, slog := range slogLoggers {
		go slog.Error(msg, args...)
	}
}
