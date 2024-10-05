package errs

import (
	"errors"
	"fmt"
)

// errorString - improved error structure that stores a message and the original error.
type errorString struct {
	message string // Detailed error message.
	origErr string // Original error.
}

// New returns a new error that includes a message and the original error.
//
// The function takes a single parameter:
// - message: A string that represents the detailed error message.
//
// The function returns an error that contains the provided message and the original error.
// If the original error is not provided, the message is used as the original error as well.
//
// Example:
//
//	err := New("Failed to connect to the database")
//	fmt.Println(err.Error()) // Output: Failed to connect to the database
func New(message string) error {
	return &errorString{
		message: message,
		origErr: message,
	}
}

// NewF creates a new error with a formatted message and the original error.//+
// It uses the fmt.Sprintf function to format the message using the provided format and arguments.//+
// The function returns a new errorString instance with the formatted message and the original error.//+
//
// Parameters:
// - format: A string that specifies the format of the error message.//+
// - a: A variadic parameter that accepts any number of arguments to be used in the format string.//+
//
// Returns:
// - An error that contains the formatted message and the original error.//+
func NewF(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return &errorString{
		message: msg,
		origErr: msg,
	}
}

// Error implements the error interface, returning the original error message.
// If the errorString instance is nil, it returns an empty string.
//
// Parameters:
// - e: A pointer to the errorString instance.
//
// Returns:
// - A string representing the original error message.
func (e *errorString) Error() string {
	if e == nil {
		return ""
	}

	return e.origErr
}

// Is checks if the target error is equal to the given error.
// It returns true if they are the same or if the target error matches the original error.
//
// Parameters:
// - err: The error to be checked.
// - target: The error to compare against.
//
// Returns:
//   - A boolean indicating whether the target error is equal to the given error or matches the original error.
//     If the target error is nil, it returns true if the given error is also nil.
func Is(err, target error) bool {
	if target == nil {
		return err == nil
	}
	return errors.Is(err, target)
}

// IsNil checks if the provided error is nil.
//
// The function takes a single parameter:
// - err: An error to be checked. If this parameter is nil, the function will return true.
//
// The function returns a boolean value:
// - true: If the provided error is nil.
// - false: If the provided error is not nil.
func IsNil(err error) bool {
	return err == nil
}
