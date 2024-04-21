// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors_test

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/jlourenc/xgo/xerrors"
)

func TestJoin(t *testing.T) {
	err0 := xerrors.New("err0")
	err1 := xerrors.New("err1")

	testCases := []struct {
		name     string
		errs     []error
		expected []error
	}{
		{
			name:     "join nothing",
			errs:     nil,
			expected: nil,
		},
		{
			name:     "join nil error",
			errs:     []error{nil},
			expected: nil,
		},
		{
			name:     "join nil errors",
			errs:     []error{nil, nil},
			expected: nil,
		},
		{
			name:     "join error",
			errs:     []error{err0},
			expected: []error{err0},
		},
		{
			name:     "join errors",
			errs:     []error{err0, err1},
			expected: []error{err0, err1},
		},
		{
			name:     "join errors including nil",
			errs:     []error{nil, err0, nil, err1},
			expected: []error{err0, err1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := xerrors.Join(tc.errs...)

			if tc.expected == nil && got != nil {
				t.Errorf("expected no error, got %s", got)
				return
			}

			if tc.expected != nil && got == nil {
				t.Errorf("expected %q, got no error", tc.expected)
				return
			}

			if got != nil {
				errs := got.(interface{ Unwrap() []error }).Unwrap()
				if !slices.Equal(errs, tc.expected) {
					t.Errorf("expected %v, got %v", tc.expected, got)
				}
				if len(errs) != cap(errs) {
					t.Errorf("with %v len!=cap, len=%d, cap=%d", tc.errs, len(errs), cap(errs))
				}
			}
		})
	}
}

func TestJoinError_Error(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "single error",
			err:      xerrors.Join(xerrors.New("err0")),
			expected: "err0",
		},
		{
			name: "multiple errors",
			err: xerrors.Join(
				xerrors.New("err0"),
				xerrors.Join(xerrors.New("err1.0"), xerrors.New("err1.1")),
				xerrors.Join(xerrors.New("err2")),
			),
			expected: "3 errors occurred:\n\t* err0\n\t* 2 errors occurred:\n\t\t* err1.0\n\t\t* err1.1\n\t* err2\n",
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

func TestJoinError_Format(t *testing.T) {
	xerrors.EnableStackTrace(false)

	testCases := []struct {
		name     string
		err      error
		format   string
		expected string
	}{
		{
			name: "default",
			err: xerrors.Join(
				xerrors.New("error message 0"),
				xerrors.Wrap(xerrors.Join(xerrors.New("error message 1"), xerrors.New("error message 2")), "wrapped"),
				xerrors.Join(xerrors.New("error message 3")),
			),
			format: "%v",
			expected: strings.Join([]string{
				`3 errors occurred:\n\t`,
				`\* error message 0\n\t`,
				`\* wrapped: 2 errors occurred:\n\t\t\* error message 1\n\t\t\* error message 2\n\t`,
				`\* error message 3\n`,
			}, ""),
		},
		{
			name: "default plus extra with stack trace disabled",
			err: xerrors.Join(
				xerrors.New("error message 0"),
				xerrors.Wrap(xerrors.Join(xerrors.New("error message 1"), xerrors.New("error message 2")), "wrapped"),
				xerrors.Join(xerrors.New("error message 3")),
			),
			format: "%+v",
			expected: strings.Join([]string{
				`3 errors occurred:\n\t`,
				`\* error message 0\n\t`,
				`\* wrapped: 2 errors occurred:\n\t\t\* error message 1\n\t\t\* error message 2\n\t`,
				`\* error message 3\n`,
			}, ""),
		},
		{
			name: "Go-syntax representation of the value with stack trace disabled",
			err: xerrors.Join(
				xerrors.New("error message 0"),
				xerrors.Wrap(xerrors.Join(xerrors.New("error message 1"), xerrors.New("error message 2")), "wrapped"),
				xerrors.Join(xerrors.New("error message 3")),
			),
			format:   "%#v",
			expected: `\*xerrors\.joinError\{errs:\(\[\]error\)\(0x[a-f0-9]+\)\}`,
		},
		{
			name: "string",
			err: xerrors.Join(
				xerrors.New("error message 0"),
				xerrors.Wrap(xerrors.Join(xerrors.New("error message 1"), xerrors.New("error message 2")), "wrapped"),
				xerrors.Join(xerrors.New("error message 3")),
			),
			format: "%s",
			expected: strings.Join([]string{
				`3 errors occurred:\n\t`,
				`\* error message 0\n\t`,
				`\* wrapped: 2 errors occurred:\n\t\t\* error message 1\n\t\t\* error message 2\n\t`,
				`\* error message 3\n`,
			}, ""),
		},
		{
			name: "double-quote string",
			err: xerrors.Join(
				xerrors.New("error message 0"),
				xerrors.Wrap(xerrors.Join(xerrors.New("error message 1"), xerrors.New("error message 2")), "wrapped"),
				xerrors.Join(xerrors.New("error message 3")),
			),
			format: "%q",
			expected: strings.Join([]string{
				`\"3 errors occurred:\\n\\t`,
				`\* error message 0\\n\\t`,
				`\* wrapped: 2 errors occurred:\\n\\t\\t\* error message 1\\n\\t\\t\* error message 2\\n\\t`,
				`\* error message 3\\n\"`,
			}, ""),
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

func TestJoinError_StackTrace(t *testing.T) {
	testCases := []struct {
		name             string
		errs             []error
		enableStackTrace bool
		expectedSize     int
	}{
		{
			name:         "stack error",
			errs:         []error{&stackError{}, &unstackError{}},
			expectedSize: 4,
		},
		{
			name:         "unstack error with stack trace disabled",
			errs:         []error{&unstackError{}, &stackError{}},
			expectedSize: 0,
		},
		{
			name:             "unstack error with stack trace enabled",
			errs:             []error{&unstackError{}, &stackError{}},
			enableStackTrace: true,
			expectedSize:     3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			xerrors.EnableStackTrace(tc.enableStackTrace)
			defer xerrors.EnableStackTrace(false)

			got := xerrors.Join(tc.errs...).(interface{ StackTrace() xerrors.StackTrace }).StackTrace()

			if len(got) != tc.expectedSize {
				t.Errorf("expected stack trace of size %d, got %v", tc.expectedSize, got)
			}
		})
	}
}
