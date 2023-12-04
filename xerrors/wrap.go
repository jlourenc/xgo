// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors

import (
	"errors"
	"fmt"
)

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(any) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// An error type might provide an As method so it can be treated as if it were a
// different error type.
//
// As panics if target is not a non-nil pointer to either a type that implements
// error, or to any interface type.
//
// It is the drop-in replacement for errors.As.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
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
// an example in the standard library.
//
// It is the drop-in replacement for errors.Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// Is is the drop-in replacement for errors.Unwrap.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Wrap returns an error annotating err with the supplied message and a stack trace,
// if enabled and err does not already contain one, at the point Wrap is called.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(StackTracer); !ok {
		err = &withStack{
			error: err,
			stack: callers(),
		}
	}

	return &withWrap{
		err: err,
		msg: message,
	}
}

// Wrapf returns an error annotating err with the format specifier and a stack trace,
// if enabled and err does not already contain one, at the point Wrap is called.
// If err is nil, Wrap returns nil.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(StackTracer); !ok {
		err = &withStack{
			error: err,
			stack: callers(),
		}
	}

	return &withWrap{
		err: err,
		msg: fmt.Sprintf(format, args...),
	}
}

type withWrap struct {
	err error
	msg string
}

// Error makes withWrap implement the error interface.
func (e *withWrap) Error() string { return e.msg + ": " + e.err.Error() }

// Format makes withWrap implement the fmt.Formatter interface.
func (e *withWrap) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s: %+v", e.msg, e.Unwrap())
			return
		}
		if s.Flag('#') {
			fmt.Fprintf(s, "%T{msg:%q, err:(%T)(%p)}", e, e.Error(), e.err, &e.err)
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	default:
		fmt.Fprintf(s, "&{%%!%c(%T=%s) %%!%c(%#v)}", verb, e.msg, e.Error(), verb, e.err)
	}
}

// StackTrace makes withWrap implement the StackTracer interface.
func (e *withWrap) StackTrace() StackTrace {
	return e.err.(StackTracer).StackTrace()
}

// Unwrap makes withWrap implement the errors.Unwrapper interface.
func (e *withWrap) Unwrap() error { return e.err }
