// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xerrors_test

import (
	"testing"

	"github.com/jlourenc/xgo/xerrors"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "empty message",
			message:  "",
			expected: "",
		},
		{
			name:     "non-empty message",
			message:  "a non-empty message",
			expected: "a non-empty message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := xerrors.New(tc.message)

			if tc.expected != got.Error() {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestNewf(t *testing.T) {
	testCases := []struct {
		name     string
		format   string
		args     []any
		expected string
	}{
		{
			name:     "empty format",
			format:   "",
			args:     nil,
			expected: "",
		},
		{
			name:     "non-empty format",
			format:   "a %s message",
			args:     []any{"non-empty"},
			expected: "a non-empty message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := xerrors.Newf(tc.format, tc.args...)

			if tc.expected != got.Error() {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}
