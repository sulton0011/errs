# errs

`errs` is a Go package that provides enhanced error handling and logging utilities. It offers a custom error type, logging capabilities, and convenient methods for error wrapping and unwrapping.

## Installation

You can install the package using:

```bash
go get github.com/sulton0011/errs
```

## Features

- Enhanced error structure for better error handling.
- Custom logging levels and formats (JSON, text, and file).
- Functions to wrap and log errors with context.
- Methods to check for specific error types.

## Usage

### Error Handling

#### Custom Error Type

The package provides a custom error type `errorString` to encapsulate error messages and original errors.

```go
package main

import (
	"fmt"
	"github.com/sulton0011/errs"
)

func main() {
	err := errs.New("an error occurred")
	fmt.Println(err.Error())
}
```

### Logging

#### Initialize Loggers

You can set the logging type by calling `SetLogTypes`:

```go
errs.SetLogTypes(errs.LogTypeJSON, errs.LogTypeText)
```

#### Log Errors

To log an error asynchronously, use the `Log` function:

```go
func doSomething() error {
	return errs.New("something went wrong")
}

func main() {
	err := doSomething()
	errs.Log(err, "Request data", "Additional context")
}
```

### Wrapping Errors

Use `Wrap` and `Wrapf` to add context to existing errors:

```go
func doAnotherThing() error {
	err := doSomething()
	if !errs.IsNil(err) {
		return errs.Wrap(err, "failed in doAnotherThing")
	}
	return nil
}
```

### Checking Errors

You can use the `Is` and `As` functions to check for specific error types.

### Example

Here’s a complete example demonstrating error handling, logging, and wrapping:

```go
package main

import (
	"fmt"
	"github.com/sulton0011/errs"
)

func main() {
	errs.SetLogTypes(errs.LogTypeJSON)

	err := doSomething()
	if !errs.IsNil(err) {
		errs.Log(err, "Request data", "Additional context")
	}
}

func doSomething() error {
	return errs.New("an error occurred")
}
```

## Functions

### New

```go
func New(message string) error
```

Creates a new error with a specified message.

### SetLogTypes

```go
func SetLogTypes(types ...logType)
```

Configures the logging types (e.g., JSON, text, file).

### SetLogFile

```go
func SetLogFile(filePath string) error
```

Sets the log file path and configures the file logger.

### Wrap

```go
func Wrap(err error, args ...interface{}) error
```

Adds context to an existing error by wrapping it with additional messages.

### Wrapf

```go
func Wrapf(err error, format string, args ...interface{}) error
```

Formats a message and wraps it around an existing error.

### Log

```go
func Log(err error, req interface{}, msgs ...interface{})
```

Logs an error asynchronously with additional context messages.

### Is

```go
func Is(err error, target error) bool
```

Reports whether any error in err's chain matches the target error.

### Unwrap

```go
func Unwrap(err error) string
```

Returns the original error message from a wrapped error.

### IsNil

```go
func IsNil(err error) bool
```

Checks if an error is nil.