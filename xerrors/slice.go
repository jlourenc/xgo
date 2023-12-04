// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors

import (
	"fmt"
	"strconv"
	"strings"
)

// Append is a helper function that appends errors into a single error to group
// multiple errors. Any nil error within errs is ignored. If err is not a grouped
// error then it will be turned into one.
func Append(err error, errs ...error) error {
	sliceErr, ok := err.(*withSlice)
	if !ok {
		errs = append([]error{err}, errs...)
		sliceErr = &withSlice{
			errs: make([]error, 0, len(errs)),
		}
	}

	for _, err = range errs {
		if err == nil {
			continue
		}

		if _, ok := err.(StackTracer); !ok {
			err = &withStack{
				error: err,
				stack: callers(),
			}
		}

		sliceErr.errs = append(sliceErr.errs, err)
	}

	if len(sliceErr.errs) == 0 {
		return nil
	}

	return sliceErr
}

type withSlice struct {
	errs []error
}

// Error makes withSlice implement the error interface.
func (e *withSlice) Error() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(len(e.errs)))
	if len(e.errs) > 1 {
		sb.WriteString(" errors")
	} else {
		sb.WriteString(" error")
	}
	sb.WriteString(" occurred:\n")

	for _, err := range e.errs {
		lines := strings.Split(strings.Trim(err.Error(), "\n"), "\n")
		sb.WriteString("\t* ")
		sb.WriteString(lines[0])
		sb.WriteString("\n")
		for _, line := range lines[1:] {
			sb.WriteString("\t")
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// Format makes withSlice implement the fmt.Formatter interface.
func (e *withSlice) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprint(s, strconv.Itoa(len(e.errs)))
			if len(e.errs) > 1 {
				fmt.Fprint(s, " errors")
			} else {
				fmt.Fprint(s, " error")
			}
			fmt.Fprint(s, " occurred:\n")

			for _, err := range e.errs {
				lines := strings.Split(strings.Trim(fmt.Sprintf("%+v", err), "\n"), "\n")
				fmt.Fprint(s, "\t* ")
				fmt.Fprint(s, lines[0])
				fmt.Fprint(s, "\n")
				for _, line := range lines[1:] {
					fmt.Fprint(s, "\t")
					fmt.Fprint(s, line)
					fmt.Fprint(s, "\n")
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

// StackTrace makes withSlice implement the StackTracer interface.
func (e *withSlice) StackTrace() StackTrace {
	return e.errs[0].(StackTracer).StackTrace()
}

// Unwrap makes withSlice implement the Unwrapper interface.
func (e *withSlice) Unwrap() error {
	if len(e.errs) == 1 {
		return e.errs[0]
	}

	// depth-first unwrapping
	errChain := make(chain, 0, len(e.errs))
	for _, err := range e.errs {
		for uerr := err; uerr != nil; uerr = Unwrap(uerr) {
			errChain = append(errChain, uerr)
		}
	}

	return errChain
}

type chain []error

// As implements errors.As by attempting to map to the current value.
func (e chain) As(target any) bool {
	return As(e[0], target)
}

// Error makes chain implement the error interface.
func (e chain) Error() string {
	return e[0].Error()
}

// Is implements errors.Is by comparing the current value directly.
func (e chain) Is(target error) bool {
	return Is(e[0], target)
}

// StackTrace makes chain implement the StackTracer interface.
func (e chain) StackTrace() StackTrace {
	if st, ok := e[0].(StackTracer); ok {
		return st.StackTrace()
	}
	return nil
}

// Unwrap makes chain implement the Unwrapper interface.
func (e chain) Unwrap() error {
	if len(e) == 1 {
		return nil
	}
	return e[1:]
}
