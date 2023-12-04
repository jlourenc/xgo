// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xerrors implements functions to manipulate errors.
// It is meant as a drop-in replacement for the Go standard library error package
// with additional functionalities inspired by github.com/pkg/errors and
// github.com/hashicorp/go-multierror packages.
package xerrors

import (
	"errors"
	"fmt"
)

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
// New also records a stack trace at the point it is called if enabled.
//
// It is the drop-in replacement for errors.New.
func New(message string) error {
	return &withStack{
		error: errors.New(message),
		stack: callers(),
	}
}

// Newf formats according to a format specifier and returns the string as a value that satisfies error.
// New also records a stack trace at the point it is called if enabled.
//
// It is the equivalent of fmt.Errorf.
func Newf(format string, args ...any) error {
	return &withStack{
		error: fmt.Errorf(format, args...),
		stack: callers(),
	}
}
