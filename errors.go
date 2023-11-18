// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// because the former will succeed if err wraps an *fs.PathError.
package errs

import (
	"fmt"
	"log/slog"
	"os"
)

// type Error struct {
// 	Err error // error response
// 	errorLog error // error logs
// }

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text, text}
}

type errorString struct {
	errinfo string
	errlog  string
}

func (e *errorString) Error() string {
	return e.errinfo
}

func (e *errorString) ErrorLog() *errorString {
	return e
}

func Wrap(e *error, msgs ...any) error {
	if *e == nil {
		return nil
	}

	err := errorlog((*e))

	var nonNilMsgs []interface{}
	for _, msg := range msgs {
		if msg != nil {
			nonNilMsgs = append(nonNilMsgs, "--->", msg)
		}
	}

	if len(nonNilMsgs) != 0 {
		err.errlog = fmt.Sprintln(nonNilMsgs[1:]...)
	}

	(*e) = err
	return (*e)
}

func WrapLog(err *error, req interface{}, msgs ...interface{}) {
	if *err == nil {
		return
	}

	var nonNilMsgs []interface{}
	for _, msg := range msgs {
		if msg != nil {
			nonNilMsgs = append(nonNilMsgs, "--->", msg)
		}
	}

	slogs := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelError,
	}))

	slogs.Error(
		fmt.Sprintln(nonNilMsgs[1:]...),
		slog.String("Error", errorlog(*err).errlog),
		slog.Any("request", req),
	)
}

func errorlog(e error) *errorString {
	err, ok := e.(interface {
		ErrorLog() *errorString
	})

	if !ok {
		if e == nil {
			return nil
		}
		return errorlog(New(e.Error()))
	}

	return err.ErrorLog()
}
