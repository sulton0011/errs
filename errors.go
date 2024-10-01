package errs

import "errors"

// errorString - improved error structure that stores a message and the original error.
type errorString struct {
	message string // Detailed error message.
	origErr string // Original error.
}

// New returns a new error that includes a message and the original error.
func New(message string) error {
	return &errorString{
		message: message,
		origErr: message,
	}
}

// Error implements the error interface, returning the original error message.
func (e *errorString) Error() string {
	return e.origErr
}

// Is checks if the target error is equal to the given error.
// It returns true if they are the same or if the target error matches the original error.
func Is(err, target error) bool {
	if target == nil {
		return err == nil
	}
	return errors.Is(err, target)
}

// IsNil checks if the provided error is nil.
// It returns true if the error is nil, false otherwise.
func IsNil(err error) bool {
	return err == nil
}