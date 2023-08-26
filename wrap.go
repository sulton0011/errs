// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errs

import (
	"errors"
	"fmt"
	"log"
)

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// Unwrap returns nil if the Unwrap method returns []error.
func Unwrap(err Error) error {
	u, ok := err.Err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

// Is reports whether any error in err's tree matches target.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly
// calling Unwrap. When err wraps multiple errors, Is examines err followed by a
// depth-first traversal of its children.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library. An Is method should only shallowly
// compare err and the target and not call Unwrap on either.
func Is(err, target Error) bool {
	return errors.Is(err.Err, target.Err)
}

// As finds the first error in err's tree that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly
// calling Unwrap. When err wraps multiple errors, As examines err followed by a
// depth-first traversal of its children.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// An error type might provide an As method so it can be treated as if it were a
// different error type.
//
// As panics if target is not a non-nil pointer to either a type that implements
// error, or to any interface type.
func As(err Error, target any) bool {
	return errors.As(err.Err, target)
}

// Wrap wraps an error with additional messages.
//
// If the given error err is not nil, Wrap creates a new error by combining
// messages from the msgs arguments with the text of the original error.
// It returns a new error with the combined data.
// If err is nil, the function returns nil.
func Wrap(err *error, msgs ...interface{}) error {
	if err == nil {
		return *err
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
		*err = New(message + "--->" + (*err).Error())
	}

	return *err
}

// WrapLog wraps an error with additional messages and performs logging.
//
// If the given error err is not nil, WrapLog creates a message by combining
// messages from the msgs arguments with the text of the original error.
// It then logs this message along with information about the request and error.
// If err is nil, the function does not take any action.
func WrapLog(err *error, req interface{}, msgs ...interface{}) {
	if *err == nil {
		return
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

	log.Print(
		message,
		" request: ", req,
		" Error: ", *err,
	)
}
