// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// because the former will succeed if err wraps an *fs.PathError.
package errs

import "fmt"

// type Error struct {
// 	Err error // error response
// 	errorLog error // error logs
// }

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) Error {
	return &errorString{text, text}
}

type Error interface {
	Error() string
	Wrap(...any) error
}

// errorString is a trivial implementation of error.
type errorString struct {
	errinfo string
	errlog  string
}

func (e *errorString) Error() string {
	return e.errinfo
}

func (e *errorString) Wrap(msgs ...any) error {
	if e == nil {
		return nil
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
		e.errlog = message + "--->" + e.errlog
	}

	return e
}
