// Copyright 2021 Jérémy Lourenço. All rights reserved.
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

type (
	stackError   struct{}
	unstackError struct{}
)

func (unstackError) Error() string        { return "unstack error" }
func (stackError) Error() string          { return "stack error" }
func (stackError) StackTrace() StackTrace { return []Frame{0, 1, 2, 3} }

func TestWithStack(t *testing.T) {
	testCases := []struct {
		name   string
		err    error
		assert func(t *testing.T, err error)
	}{
		{
			name: "nil",
			err:  nil,
			assert: func(t *testing.T, err error) {
				if err != nil {
					t.Errorf("expected nil, got %#v", err)
				}
			},
		},
		{
			name: "stack error",
			err:  stackError{},
			assert: func(t *testing.T, err error) {
				if errors.Unwrap(err) != nil || !errors.As(err, &stackError{}) {
					t.Errorf("expected %#v, got %#v", stackError{}, err)
				}
			},
		},
		{
			name: "unstack error",
			err:  unstackError{},
			assert: func(t *testing.T, err error) {
				if errors.Unwrap(err) == nil || !errors.As(err, &unstackError{}) {
					t.Errorf("expected %#v, got %#v", stackError{}, err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := WithStack(tc.err)

			tc.assert(t, got)
		})
	}
}

func TestWithStack_Format(t *testing.T) {
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
			expected: `^error message$`,
		},
		{
			name:     "default format plus extra with stack trace disabled",
			err:      errors.New("error message"),
			format:   "%+v",
			expected: `^error message$`,
		},
		{
			name:             "default format plus extra with stack trace enabled",
			err:              errors.New("error message"),
			enableStackTrace: true,
			format:           "%+v",
			expected:         `^error message(\n(\t)?[0-9a-zA-Z.\/_:-]+)+$`,
		},
		{
			name:     "Go-syntax representation of the value with stack trace disabled",
			err:      errors.New("error message"),
			format:   "%#v",
			expected: `^\*xerrors\.withStack\{error:\(string\)\(0x[a-f0-9]+\), stack:xerrors\.stack\(nil\)\}$`,
		},
		{
			name:             "Go-syntax representation of the value with stack trace enabled",
			err:              errors.New("error message"),
			enableStackTrace: true,
			format:           "%#v",
			expected:         `^\*xerrors\.withStack\{error:\(string\)\(0x[a-f0-9]+\), stack:xerrors\.stack\[([0-9]+[ ]?){3}\]\}$`,
		},
		{
			name:     "string format",
			err:      errors.New("error message"),
			format:   "%s",
			expected: `^error message$`,
		},
		{
			name:     "double-quoted string format",
			err:      errors.New("error message"),
			format:   "%q",
			expected: `^"error message"$`,
		},
		{
			name:     "Go-syntax representation of the type of the value",
			err:      errors.New("error message"),
			format:   "%T",
			expected: `^\*xerrors\.withStack$`,
		},
		{
			name:             "unsupported format",
			err:              errors.New("error message"),
			enableStackTrace: true,
			format:           "%t",
			expected:         `^\&\{\%\!t\(\&errors\.errorString\{s:"error message"\}\) \[(\%\!t\(uintptr=[0-9]+\)[ ]?){3}\]\}$`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			EnableStackTrace(tc.enableStackTrace)
			defer EnableStackTrace(false)

			got := fmt.Sprintf(tc.format, WithStack(tc.err))

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

func TestWithStack_Unwrap(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "unwrap",
			err:      WithStack(&unstackError{}),
			expected: &unstackError{},
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
