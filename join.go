package errs

import (
	"fmt"
	"strings"
)

// JoinError concatenates multiple error messages into a single error object using the specified sep.
// It skips nil errors and ensures a clean, unified error message string.
//
// Parameters:
//   - sep: A string representing the sep to be used between the error messages.
//   - errors: Variadic arguments representing the errors to be concatenated.
//
// Returns:
//   - An error object containing the concatenated error messages, separated by the specified sep.
//     If all errors are nil, it returns nil.
//     The error object also contains the original error messages in its error string representation.
func Join(sep string, errors ...error) error {
	var origErr strings.Builder
	var message strings.Builder
	for _, err := range errors {
		if err != nil {
			if origErr.Len() > 0 {
				origErr.WriteString(sep) // Append sep between error messages.
				message.WriteString(sep) // Append sep between error messages.
			}
			origErr.WriteString(err.Error()) // Append error message to builder.
			message.WriteString(Unwrap(err)) // Append error message to builder.
		}
	}

	if origErr.Len() == 0 {
		return nil
	}
	
	return &errorString{
		origErr: origErr.String(),
		message: message.String(),
	}
}

// JoinMessages concatenates multiple messages using the specified sep.
// It skips nil values and ensures a clean, unified message string.
//
// Parameters:
//   - sep: A string representing the sep to be used between the messages.
//   - messages: Variadic arguments representing the messages to be concatenated.
//
// Returns:
//   - A string representing the concatenated messages, separated by the specified sep.
//     If all messages are nil, an empty string is returned.
func JoinMsg(sep string, a ...any) string {
	var builder strings.Builder
	for _, message := range a {
		if message != nil {
			if builder.Len() > 0 {
				builder.WriteString(sep) // Append sep between messages.
			}
			builder.WriteString(fmt.Sprint(message)) // Convert and append message to builder.
		}
	}
	return builder.String()
}
