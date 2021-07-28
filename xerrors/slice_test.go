// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors_test

import (
	"fmt"
	"regexp"
	"testing"

	. "github.com/jlourenc/xgo/xerrors"
)

var isErr = New("is error")

func TestAppend(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		errs     []error
		expected string // empty string means no error
	}{
		{
			name:     "append nil error array to nil error",
			err:      nil,
			errs:     nil,
			expected: "",
		},
		{
			name:     "append empty error array to nil error",
			err:      nil,
			errs:     []error{},
			expected: "",
		},
		{
			name:     "append nil errors to nil error",
			err:      nil,
			errs:     []error{nil, nil},
			expected: "",
		},
		{
			name:     "append nil error array to error",
			err:      &stackError{},
			errs:     nil,
			expected: "1 error occurred:\n\t* stack error\n",
		},
		{
			name:     "append empty error array to error",
			err:      &stackError{},
			errs:     []error{},
			expected: "1 error occurred:\n\t* stack error\n",
		},
		{
			name:     "append nil errors to error",
			err:      &unstackError{},
			errs:     []error{nil, nil},
			expected: "1 error occurred:\n\t* unstack error\n",
		},
		{
			name:     "append errors to nil error",
			err:      nil,
			errs:     []error{&stackError{}, &unstackError{}},
			expected: "2 errors occurred:\n\t* stack error\n\t* unstack error\n",
		},
		{
			name:     "append errors to a grouped error",
			err:      Append(&stackError{}),
			errs:     []error{&stackError{}, &unstackError{}},
			expected: "3 errors occurred:\n\t* stack error\n\t* stack error\n\t* unstack error\n",
		},
		{
			name:     "append grouped errors to a group error",
			err:      Append(&stackError{}),
			errs:     []error{Append(&stackError{}), Append(&unstackError{})},
			expected: "3 errors occurred:\n\t* stack error\n\t* 1 error occurred:\n\t\t* stack error\n\t* 1 error occurred:\n\t\t* unstack error\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Append(tc.err, tc.errs...)

			if tc.expected == "" && got != nil {
				t.Errorf("expected no error, got %s", got)
			} else if tc.expected != "" {
				if got == nil {
					t.Errorf("expected %q, got no error", tc.expected)
				} else if tc.expected != got.Error() {
					t.Errorf("expected %q, got %q", tc.expected, got)
				}
			}
		})
	}
}

func TestWithSlice_As(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "error type is part of a simple grouped error",
			err:      Append(&unstackError{}),
			expected: true,
		},
		{
			name:     "error type is not part of a simple grouped error",
			err:      Append(&stackError{}),
			expected: false,
		},
		{
			name:     "error type is part of a multiple grouped error",
			err:      Append(Append(New("error message 0"), Wrap(&unstackError{}, "wrapped")), New("error message 1")),
			expected: true,
		},
		{
			name:     "error type is not part of a multiple grouped error",
			err:      Append(Append(New("error message 0"), Wrap(&stackError{}, "wrapped")), New("error message 1")),
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var uerr *unstackError
			got := As(tc.err, &uerr)

			if tc.expected != got {
				t.Errorf("expected %t, got %t", tc.expected, got)
			}
		})
	}
}

func TestWithSlice_Error(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "error",
			err:      Append(New("error message 0"), Wrap(Append(New("error message 1"), New("error message 2")), "wrapped"), Append(New("error message 3"))),
			expected: "3 errors occurred:\n\t* error message 0\n\t* wrapped: 2 errors occurred:\n\t\t* error message 1\n\t\t* error message 2\n\t* 1 error occurred:\n\t\t* error message 3\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.err.Error()

			if tc.expected != got {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestWithSlice_Format(t *testing.T) {
	EnableStackTrace(false)

	testCases := []struct {
		name     string
		err      error
		format   string
		expected string
	}{
		{
			name:     "default",
			err:      Append(New("error message 0"), Wrap(Append(New("error message 1"), New("error message 2")), "wrapped"), Append(New("error message 3"))),
			format:   "%v",
			expected: `3 errors occurred:\n\t\* error message 0\n\t\* wrapped: 2 errors occurred:\n\t\t\* error message 1\n\t\t\* error message 2\n\t\* 1 error occurred:\n\t\t\* error message 3\n`,
		},
		{
			name:     "default plus extra with stack trace disabled",
			err:      Append(New("error message 0"), Wrap(Append(New("error message 1"), New("error message 2")), "wrapped"), Append(New("error message 3"))),
			format:   "%+v",
			expected: `3 errors occurred:\n\t\* error message 0\n\t\* wrapped: 2 errors occurred:\n\t\t\* error message 1\n\t\t\* error message 2\n\t\* 1 error occurred:\n\t\t\* error message 3\n`,
		},
		{
			name:     "Go-syntax representation of the value with stack trace disabled",
			err:      Append(New("error message 0"), Wrap(Append(New("error message 1"), New("error message 2")), "wrapped"), Append(New("error message 3"))),
			format:   "%#v",
			expected: `\*xerrors\.withSlice\{errs:\(\[\]error\)\(0x[a-f0-9]+\)\}`,
		},
		{
			name:     "string",
			err:      Append(New("error message 0"), Wrap(Append(New("error message 1"), New("error message 2")), "wrapped"), Append(New("error message 3"))),
			format:   "%s",
			expected: `3 errors occurred:\n\t\* error message 0\n\t\* wrapped: 2 errors occurred:\n\t\t\* error message 1\n\t\t\* error message 2\n\t\* 1 error occurred:\n\t\t\* error message 3\n`,
		},
		{
			name:     "double-quote string",
			err:      Append(New("error message 0"), Wrap(Append(New("error message 1"), New("error message 2")), "wrapped"), Append(New("error message 3"))),
			format:   "%q",
			expected: `\"3 errors occurred:\\n\\t\* error message 0\\n\\t\* wrapped: 2 errors occurred:\\n\\t\\t\* error message 1\\n\\t\\t\* error message 2\\n\\t\* 1 error occurred:\\n\\t\\t\* error message 3\\n\"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fmt.Sprintf(tc.format, tc.err)

			re, err := regexp.Compile(tc.expected)
			if err != nil {
				t.Fatalf("invalid regex: %s", tc.expected)
			}
			if !re.MatchString(got) {
				t.Errorf("expected pattern %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestWithSlice_Is(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "error type part of a simple grouped error",
			err:      Append(New("error message 0"), Wrap(isErr, "wrapped"), New("error message 1")),
			expected: true,
		},
		{
			name:     "error type not part of a simple grouped error",
			err:      Append(New("error message 0"), Wrap(&stackError{}, "wrapped"), New("error message 1")),
			expected: false,
		},
		{
			name:     "error type part of a multiple grouped error",
			err:      Append(Append(New("error message 0"), Wrap(isErr, "wrapped")), New("error message 1")),
			expected: true,
		},
		{
			name:     "error type not part of a multiple grouped error",
			err:      Append(Append(New("error message 0"), Wrap(&stackError{}, "wrapped")), New("error message 1")),
			expected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Is(tc.err, isErr)

			if tc.expected != got {
				t.Errorf("expected %t, got %t", tc.expected, got)
			}
		})
	}
}

func TestWithSlice_StackTrace(t *testing.T) {
	testCases := []struct {
		name             string
		err              error
		enableStackTrace bool
		expectedSize     int
	}{
		{
			name:         "stack error",
			err:          &stackError{},
			expectedSize: 4,
		},
		{
			name:         "unstack error with stack trace disabled",
			err:          &unstackError{},
			expectedSize: 0,
		},
		{
			name:             "unstack error with stack trace enabled",
			err:              &unstackError{},
			enableStackTrace: true,
			expectedSize:     3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			EnableStackTrace(tc.enableStackTrace)
			defer EnableStackTrace(false)

			got := Append(tc.err).(interface{ StackTrace() StackTrace }).StackTrace()

			if len(got) != tc.expectedSize {
				t.Errorf("expected stack trace of size %d, got %v", tc.expectedSize, got)
			}
		})
	}
}

func TestChain_Error(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "error 0",
			err:      Unwrap(Append(New("error message 0"), New("error message 1"), New("error message 3"))),
			expected: "error message 0",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.err.Error()

			if tc.expected != got {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestChain_StackTrace(t *testing.T) {
	testCases := []struct {
		name             string
		err              error
		enableStackTrace bool
		expectedSize     int
	}{
		{
			name:         "stack error",
			err:          Unwrap(Append(&stackError{}, New("error message 0"))),
			expectedSize: 4,
		},
		{
			name:         "unstack error with stack trace disabled",
			err:          Unwrap(Append(&unstackError{}, New("error message 0"))),
			expectedSize: 0,
		},
		{
			name:         "unwrapped error with no stack trace",
			err:          Unwrap(Unwrap(Append(&unstackError{}, New("error message 0")))),
			expectedSize: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			EnableStackTrace(tc.enableStackTrace)
			defer EnableStackTrace(false)

			got := tc.err.(interface{ StackTrace() StackTrace }).StackTrace()

			if len(got) != tc.expectedSize {
				t.Errorf("expected stack trace of size %d, got %v", tc.expectedSize, got)
			}
		})
	}
}
