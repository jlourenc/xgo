// Copyright 2022 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	. "github.com/jlourenc/xgo/xerrors"
)

func TestAs(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		target   interface{}
		expected bool
	}{
		{
			name:     "matches",
			err:      fmt.Errorf("wrapped error: %w", unstackError{}),
			target:   &unstackError{},
			expected: true,
		},
		{
			name:     "does not match",
			err:      fmt.Errorf("wrapped error: %w", stackError{}),
			target:   &unstackError{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := As(tc.err, tc.target)

			if tc.expected != got {
				t.Errorf("expected %t, got %t", tc.expected, got)
			}
		})
	}
}

func TestIs(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		target   error
		expected bool
	}{
		{
			name:     "matches",
			err:      fmt.Errorf("wrapped error: %w", unstackError{}),
			target:   unstackError{},
			expected: true,
		},
		{
			name:     "does not match",
			err:      fmt.Errorf("wrapped error: %w", stackError{}),
			target:   unstackError{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Is(tc.err, tc.target)

			if tc.expected != got {
				t.Errorf("expected %t, got %t", tc.expected, got)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "simple error",
			err:      unstackError{},
			expected: nil,
		},
		{
			name:     "wrapped error",
			err:      fmt.Errorf("wrapped error: %w", unstackError{}),
			expected: unstackError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Unwrap(tc.err)

			if tc.expected != got {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	testCases := []struct {
		name   string
		err    error
		assert func(t *testing.T, err error)
	}{
		{
			name: "nil",
			err:  nil,
			assert: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Errorf("expected nil, got %#v", err)
				}
			},
		},
		{
			name: "stack error",
			err:  stackError{},
			assert: func(t *testing.T, err error) {
				t.Helper()
				if errors.Unwrap(err) == nil || !errors.As(err, &stackError{}) || err.Error() != "wrap: stack error" {
					t.Errorf("expected %#v, got %#v", stackError{}, err)
				}
			},
		},
		{
			name: "unstack error",
			err:  unstackError{},
			assert: func(t *testing.T, err error) {
				t.Helper()
				if errors.Unwrap(err) == nil || !errors.As(err, &unstackError{}) || err.Error() != "wrap: unstack error" {
					t.Errorf("expected %#v, got %#v", stackError{}, err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Wrap(tc.err, "wrap")

			tc.assert(t, got)
		})
	}
}

func TestWrapf(t *testing.T) {
	testCases := []struct {
		name   string
		err    error
		assert func(t *testing.T, err error)
	}{
		{
			name: "nil",
			err:  nil,
			assert: func(t *testing.T, err error) {
				t.Helper()
				if err != nil {
					t.Errorf("expected nil, got %#v", err)
				}
			},
		},
		{
			name: "stack error",
			err:  stackError{},
			assert: func(t *testing.T, err error) {
				t.Helper()
				if errors.Unwrap(err) == nil || !errors.As(err, &stackError{}) || err.Error() != "wrap: stack error" {
					t.Errorf("expected %#v, got %#v", stackError{}, err)
				}
			},
		},
		{
			name: "unstack error",
			err:  unstackError{},
			assert: func(t *testing.T, err error) {
				t.Helper()
				if errors.Unwrap(err) == nil || !errors.As(err, &unstackError{}) || err.Error() != "wrap: unstack error" {
					t.Errorf("expected %#v, got %#v", stackError{}, err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Wrapf(tc.err, "%s", "wrap")

			tc.assert(t, got)
		})
	}
}

func TestWithWrap_Error(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "error",
			err:      Wrap(&stackError{}, "wrap"),
			expected: "wrap: stack error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.err.Error()

			if tc.expected != got {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestWithWrap_Format(t *testing.T) {
	testCases := []struct {
		name             string
		err              error
		enableStackTrace bool
		format           string
		expected         string
	}{
		{
			name:     "default format",
			err:      errors.New("error message"),
			format:   "%v",
			expected: `wrapped: error message`,
		},
		{
			name:     "default format plus extra with stack trace disabled",
			err:      Wrap(errors.New("error message"), "wrapped 0"),
			format:   "%+v",
			expected: `wrapped: wrapped 0: error message`,
		},
		{
			name:             "default format plus extra with stack trace enabled",
			err:              errors.New("error message"),
			enableStackTrace: true,
			format:           "%+v",
			expected:         `^wrapped: error message(\n(\t)?[0-9a-zA-Z.\/_:-]+)+$`,
		},
		{
			name:     "Go-syntax representation of the value with stack trace disabled",
			err:      errors.New("error message"),
			format:   "%#v",
			expected: `\*xerrors\.withWrap\{msg:"wrapped: error message", err:\(\*xerrors\.withStack\)\(0x[a-f0-9]+\)\}`,
		},
		{
			name:             "Go-syntax representation of the value with stack trace enabled",
			err:              errors.New("error message"),
			enableStackTrace: true,
			format:           "%#v",
			expected:         `\*xerrors\.withWrap\{msg:"wrapped: error message", err:\(\*xerrors\.withStack\)\(0x[a-f0-9]+\)\}`,
		},
		{
			name:     "string format",
			err:      errors.New("error message"),
			format:   "%s",
			expected: `wrapped: error message`,
		},
		{
			name:     "double-quoted string format",
			err:      errors.New("error message"),
			format:   "%q",
			expected: `"wrapped: error message"`,
		},
		{
			name:     "Go-syntax representation of the type of the value",
			err:      errors.New("error message"),
			format:   "%T",
			expected: `\*xerrors\.withWrap`,
		},
		{
			name:             "unsupported format",
			err:              errors.New("error message"),
			enableStackTrace: true,
			format:           "%t",
			expected:         `\&\{\%\!t\(string=wrapped: error message\) \%\!t\(\*xerrors\.withStack\{error:\(string\)\(0x[a-f0-9]+\), stack:xerrors\.stack\[([0-9]+[ ]?){3}\]\}\)\}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			EnableStackTrace(tc.enableStackTrace)
			defer EnableStackTrace(false)

			got := fmt.Sprintf(tc.format, Wrap(tc.err, "wrapped"))

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

func TestWithWrap_StackTrace(t *testing.T) {
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

			got := Wrap(tc.err, "wrapped").(interface{ StackTrace() StackTrace }).StackTrace()

			if len(got) != tc.expectedSize {
				t.Errorf("expected stack trace of size %d, got %v", tc.expectedSize, got)
			}
		})
	}
}

func TestWithWrap_Unwrap(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "unwrap",
			err:      Wrap(&stackError{}, ""),
			expected: &stackError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.err.(interface{ Unwrap() error }).Unwrap()

			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
