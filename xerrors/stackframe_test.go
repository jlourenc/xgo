// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors_test

import (
	"fmt"
	"regexp"
	"runtime"
	"testing"

	"github.com/jlourenc/xgo/xerrors"
)

func TestFrame_Format(t *testing.T) {
	var pcs [1]uintptr
	runtime.Callers(1, pcs[:])

	testCases := []struct {
		name     string
		format   string
		frame    xerrors.Frame
		expected string
	}{
		{
			name:     "source file",
			format:   "%s",
			frame:    xerrors.Frame(pcs[0]),
			expected: `^stackframe_test.go$`,
		},
		{
			name:     "source file plus extra",
			format:   "%+s",
			frame:    xerrors.Frame(pcs[0]),
			expected: `^github\.com\/jlourenc\/xgo\/xerrors_test\.TestFrame_Format\n\t.*\/xgo\/xerrors\/stackframe_test\.go$`,
		},
		{
			name:     "source line",
			format:   "%d",
			frame:    xerrors.Frame(pcs[0]),
			expected: `^18$`,
		},
		{
			name:     "function name",
			format:   "%n",
			frame:    xerrors.Frame(pcs[0]),
			expected: `^TestFrame_Format$`,
		},
		{
			name:     "source file and line",
			format:   "%v",
			frame:    xerrors.Frame(pcs[0]),
			expected: `^stackframe_test\.go:18$`,
		},
		{
			name:     "source file and line plus extra of unknown frame",
			format:   "%+v",
			frame:    xerrors.Frame(0),
			expected: `^unknown\n\tunknown:0$`,
		},
		{
			name:     "source file and line plus extra",
			format:   "%+v",
			frame:    xerrors.Frame(pcs[0]),
			expected: `^github.com\/jlourenc\/xgo\/xerrors_test\.TestFrame_Format\n\t.*\/xgo\/xerrors\/stackframe_test\.go:18$`,
		},
		{
			name:     "unsupported format",
			format:   "%t",
			frame:    xerrors.Frame(pcs[0]),
			expected: `^$`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fmt.Sprintf(tc.format, tc.frame)

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

func TestFrame_MarshalText(t *testing.T) {
	var pcs [1]uintptr
	runtime.Callers(1, pcs[:])

	testCases := []struct {
		name     string
		frame    xerrors.Frame
		expected string
	}{
		{
			name:     "known frame",
			frame:    xerrors.Frame(pcs[0]),
			expected: `^github.com\/jlourenc\/xgo\/xerrors_test\.TestFrame_MarshalText .*\/xgo\/xerrors\/stackframe_test\.go:93$`,
		},
		{
			name:     "unknown frame",
			frame:    xerrors.Frame(0),
			expected: `^unknown$`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.frame.MarshalText()
			if err != nil {
				t.Fatalf("error is not expected")
			}

			re, err := regexp.Compile(tc.expected)
			if err != nil {
				t.Fatalf("invalid regex: %s", tc.expected)
			}
			if !re.Match(got) {
				t.Errorf("expected pattern %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestStackTrace_Format(t *testing.T) {
	var pcs [2]uintptr
	runtime.Callers(1, pcs[:])

	testCases := []struct {
		name       string
		format     string
		stackTrace xerrors.StackTrace
		expected   string
	}{
		{
			name:       "lists source files",
			format:     "%s",
			stackTrace: xerrors.StackTrace{xerrors.Frame(pcs[0]), xerrors.Frame(pcs[1])},
			expected:   `^\[stackframe_test\.go testing\.go\]$`,
		},
		{
			name:       "lists source files and line numbers",
			format:     "%v",
			stackTrace: xerrors.StackTrace{xerrors.Frame(pcs[0])},
			expected:   `^\[stackframe_test\.go:132\]$`,
		},
		{
			name:       "lists source files, line numbers and function names",
			format:     "%+v",
			stackTrace: xerrors.StackTrace{xerrors.Frame(pcs[0])},
			expected:   `^\ngithub\.com\/jlourenc\/xgo\/xerrors_test\.TestStackTrace_Format\n\t.*\/xgo\/xerrors\/stackframe_test\.go:132$`,
		},
		{
			name:       "source file and line plus extra",
			format:     "%#v",
			stackTrace: xerrors.StackTrace{xerrors.Frame(pcs[0])},
			expected:   `^\[\]xerrors\.Frame\{stackframe_test\.go:132\}$`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := fmt.Sprintf(tc.format, tc.stackTrace)

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
