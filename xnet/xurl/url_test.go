// Copyright 2022 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xurl_test

import (
	"testing"

	. "github.com/jlourenc/xgo/xnet/xurl"
)

func TestJoinBasePath(t *testing.T) {
	testCases := []struct {
		name     string
		base     string
		elems    []string
		expected string
	}{
		{
			name:     "nil elements",
			base:     "http://localhost:80",
			elems:    nil,
			expected: "http://localhost:80/",
		},
		{
			name:     "empty elements",
			base:     "http://localhost:80",
			elems:    []string{},
			expected: "http://localhost:80/",
		},
		{
			name:     "single element",
			base:     "http://localhost:80",
			elems:    []string{"segment"},
			expected: "http://localhost:80/segment",
		},
		{
			name:     "multiple elements",
			base:     "http://localhost:80",
			elems:    []string{"segment1", "segment2", "segment_to_%escape"},
			expected: "http://localhost:80/segment1/segment2/segment_to_%25escape",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := JoinBasePath(tc.base, tc.elems...)

			if tc.expected != got {
				t.Errorf("expected %v; got %v", tc.expected, got)
			}
		})
	}
}

func TestJoinPath(t *testing.T) {
	testCases := []struct {
		name     string
		elems    []string
		expected string
	}{
		{
			name:     "nil elements",
			elems:    nil,
			expected: "",
		},
		{
			name:     "empty elements",
			elems:    []string{},
			expected: "",
		},
		{
			name:     "single element",
			elems:    []string{"segment"},
			expected: "segment",
		},
		{
			name:     "multiple elements",
			elems:    []string{"segment1", "segment2", "segment_to_%escape"},
			expected: "segment1/segment2/segment_to_%25escape",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := JoinPath(tc.elems...)

			if tc.expected != got {
				t.Errorf("expected %v; got %v", tc.expected, got)
			}
		})
	}
}
