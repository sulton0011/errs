# `errs` Package

The `errs` package provides utilities for error handling, wrapping, and structured logging in Go. It enhances the standard error handling capabilities by adding additional context and structured logging options.

## Table of Contents

- [Installation](#installation)
- [Features](#features)
- [Usage](#usage)
  - [Creating Errors](#creating-errors)
  - [Wrapping Errors](#wrapping-errors)
  - [Comparing Errors](#comparing-errors)
  - [Joining Errors](#joining-errors)
  - [Logging](#logging)
- [Configuration](#configuration)
- [Functions](#functions)
  - [New](#new)
  - [NewF](#newf)
  - [SetLogTypes](#setlogtypes)
  - [SetLogFile](#setlogfile)
  - [Wrap](#wrap)
  - [WrapF](#wrapf)
  - [Unwrap](#unwrap)
  - [UnwrapE](#unwrape)
  - [Log](#log)
  - [Join](#join)
  - [JoinMsg](#joinmsg)
  - [Is](#is)
  - [IsNil](#isnil)

## Installation

To install the package, run:

```bash
go get github.com/sulton0011/errs
```

## Features

- **Structured Error Handling**: Create and manage errors with detailed context.
- **Error Wrapping**: Use `Wrap` and `WrapF` to add context to existing errors.
- **Custom Loggers**: Support for JSON, text, and file-based logging through `LogType` configuration.
- **Error Joining**: Combine multiple errors into a single message.
- **Structured Logging**: Log errors in a structured manner using the `slog` library with various output formats.

## Usage

### Creating Errors

You can create new errors using `New` or `NewF`:

```go
package main

import "github.com/sulton0011/errs"

func main() {
    err := errs.New("simple error")
    errF := errs.NewF("formatted error: %s", "some detail")
    
    println(err.Error())     // Output: simple error
    println(errF.Error())    // Output: formatted error: some detail
}
```

### Wrapping Errors

To add context to existing errors, use `Wrap`:

```go
package main

import (
    "github.com/sulton0011/errs"
    "fmt"
)

func main() {
    baseErr := errs.New("file not found")
    wrappedErr := errs.Wrap(baseErr, "unable to process config file")
    
    fmt.Println(wrappedErr.Error()) // Output: file not found
}
```

For formatted context, use `WrapF`:

```go
formattedErr := errs.WrapF(baseErr, "error in %s operation", "read")
fmt.Println(formattedErr.Error()) // Output: error in read operation: file not found
```

### Comparing Errors

Use `Is` to compare errors:

```go
if errs.Is(wrappedErr, baseErr) {
    fmt.Println("wrappedErr is equal to baseErr")
}
```

### Joining Errors

You can join multiple errors into a single message:

```go
err1 := errs.New("first error")
err2 := errs.New("second error")
joinedErr := errs.Join(" | ", err1, err2)

fmt.Println(joinedErr.Error()) // Output: first error | second error
```

### Logging

Configure logging for JSON, text, or file-based loggers:

```go
import "github.com/sulton0011/errs"

// Set up a JSON logger
errs.SetLogTypes(errs.LogTypeJSON)

// Set up a file logger
if err := errs.SetLogFile("logs/app.log"); err != nil {
    panic(err)
}
errs.SetLogTypes(errs.LogTypeFile)
```

## Configuration

### Log Type

The `LogType` enum supports three logging formats:

- `LogTypeJSON`: JSON format logging.
- `LogTypeText`: Human-readable text format logging.
- `LogTypeFile`: Logging to a specified file.

### Setting Log Types

Configure multiple logging types with:

```go
errs.SetLogTypes(errs.LogTypeJSON, errs.LogTypeText)
```

### Setting Log File

To specify a log file path for file-based logging:

```go
if err := errs.SetLogFile("log/app.log"); err != nil {
    panic("Failed to set log file: " + err.Error())
}
```

## Functions

### New

```go
func New(message string) error
```
Creates a new error with the specified message.

### NewF

```go
func NewF(format string, a ...any) error
```
Creates a new formatted error.

### SetLogTypes

```go
func SetLogTypes(types ...logType)
```
Configures the logging types (e.g., JSON, text, file).

### SetLogFile

```go
func SetLogFile(filePath string) error
```
Sets the log file path for logging.

### Wrap

```go
func Wrap(err error, args ...any) error
```
Wraps an existing error with additional context.

### Wrapf

```go
func Wrapf(err error, format string, args ...any) error
```
Formats a message and wraps it around an existing error.

### Unwrap

```go
func Unwrap(err error) string
```
Returns the original error message from a wrapped error.

### UnwrapE

```go
func UnwrapE(err error) error
```
Returns the original error from a wrapped error.

### Log

```go
func Log(err error, req any, msgs ...any)
```
Logs an error asynchronously with additional context messages.

### Join

```go
func Join(sep string, errors ...error) error
```
Joins multiple errors into a single error.

### JoinMsg

```go
func JoinMsg(sep string, a ...any) string
```
Joins multiple arguments into a single string.

### Is

```go
func Is(err error, target error) bool
```
Checks if any error in the chain matches the target error.

### IsNil

```go
func IsNil(err error) bool
```
Checks if an error is nil.