// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
// because the former will succeed if err wraps an *fs.PathError.
package errs

import "fmt"

type Error struct {
	Err error // error response
	errorLog error // error logs
}

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func (err *Error) Error() string {
	return err.Err.Error()
}

func (err *Error) Wrap(msgs ...interface{}) *Error {
	if err.Err == nil {
		return err
	}

	if err.errorLog == nil {
		err.errorLog = err.Err
	}

	var message string
	for _, msg := range msgs {
		if msg != nil {
			if message != "" {
				message += "--->"
			}
			message += fmt.Sprint(msg)
		}
	}

	if message != "" {
		err.errorLog = New(message + "--->" + err.errorLog.Error())
	}

	return err
}
