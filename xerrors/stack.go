// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors

import (
	"fmt"
)

// WithStack annotates err with a stack trace at the point WithStack is called
// only if enabled and err does not already contain one.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(StackTracer); ok {
		return err
	}

	return &withStack{
		error: err,
		stack: callers(),
	}
}

type withStack struct {
	error
	stack
}

// Format makes withStack implement the fmt.Formatter interface.
func (e *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", e.Unwrap())
			if e.stack != nil {
				e.stack.Format(s, verb)
			}
			return
		}
		if s.Flag('#') {
			fmt.Fprintf(s, "%T{error:(%T)(%p), stack:%T%#v}", e, e.Error(), &e.error, e.stack, e.stack)
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	default:
		format := fmt.Sprintf("&{%%%%!%%c(%%#v) %%%c}", verb)
		fmt.Fprintf(s, format, verb, e.error, e.stack)
	}
}

// Unwrap makes withStack implement the errors.Unwrapper interface.
func (e *withStack) Unwrap() error { return e.error }
