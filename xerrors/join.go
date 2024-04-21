// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// Join returns an error that wraps the given errors.
// Any nil error values are discarded.
// Join returns nil if every value in errs is nil.
// The error formats as the concatenation of the strings obtained
// by calling the Error method of each element of errs, with a newline
// between each string.
//
// A non-nil error returned by Join implements the Unwrap() []error method.
//
// It is the drop-in replacement for errors.Join.
func Join(errs ...error) error {
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err == nil {
			continue
		}
		if _, ok := err.(StackTracer); !ok {
			err = &withStack{
				error: err,
				stack: callers(),
			}
		}
		e.errs = append(e.errs, err)
	}
	return e
}

type joinError struct {
	errs []error
}

// Error makes joinError implement the error interface.
func (e *joinError) Error() string {
	// Since Join returns nil if every value in errs is nil,
	// e.errs cannot be empty.
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	b := []byte(strconv.Itoa(len(e.errs)) + " errors occurred:\n")
	for _, err := range e.errs {
		b = append(b, '\t', '*', ' ')

		lines := strings.Split(strings.TrimSuffix(err.Error(), "\n"), "\n")
		b = append(b, lines[0]...)
		b = append(b, '\n')

		for _, line := range lines[1:] {
			b = append(b, '\t')
			b = append(b, line...)
			b = append(b, '\n')
		}
	}
	return unsafe.String(&b[0], len(b))
}

// Format makes joinError implement the fmt.Formatter interface.
func (e *joinError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if len(e.errs) == 1 {
				fmt.Fprint(s, e.errs[0].Error())
				return
			}

			fmt.Fprint(s, strconv.Itoa(len(e.errs)), " errors occurred:\n")
			for _, err := range e.errs {
				lines := strings.Split(strings.TrimSuffix(fmt.Sprintf("%+v", err), "\n"), "\n")
				fmt.Fprint(s, "\t* ", lines[0], "\n")
				for _, line := range lines[1:] {
					fmt.Fprint(s, "\t", line, "\n")
				}
			}
			return
		}
		if s.Flag('#') {
			fmt.Fprintf(s, "%T{errs:(%T)(%p)}", e, e.errs, &e.errs)
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// StackTrace makes joinError implement the StackTracer interface.
func (e *joinError) StackTrace() StackTrace {
	return e.errs[0].(StackTracer).StackTrace()
}

// Unwrap makes joinError implement the errors.Unwrapper interface.
func (e *joinError) Unwrap() []error {
	return e.errs
}
