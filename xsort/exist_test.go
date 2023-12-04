// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsort_test

import (
	"testing"

	"github.com/jlourenc/xgo/xsort"
)

func TestExist(t *testing.T) {
	testCases := []struct {
		name     string
		n        int
		cmp      func(int) int
		expected bool
	}{
		{
			name:     "empty",
			n:        0,
			cmp:      nil,
			expected: false,
		},
		{
			name:     "only lower values",
			n:        1e9,
			cmp:      func(i int) int { return -1 },
			expected: false,
		},
		{
			name:     "only greater values",
			n:        1e9,
			cmp:      func(i int) int { return +1 },
			expected: false,
		},
		{
			name: "item is first element - linear search",
			n:    5,
			cmp: func(i int) int {
				if i == 0 {
					return 0
				}
				return +1
			},
			expected: true,
		},
		{
			name: "item is last element - linear search",
			n:    5,
			cmp: func(i int) int {
				if i == 4 {
					return 0
				}
				return -1
			},
			expected: true,
		},
		{
			name: "item is first element - binary search",
			n:    1e9,
			cmp: func(i int) int {
				if i == 0 {
					return 0
				}
				return +1
			},
			expected: true,
		},
		{
			name: "item is last element - binary search",
			n:    1e9,
			cmp: func(i int) int {
				if i == 1e9-1 {
					return 0
				}
				return -1
			},
			expected: true,
		},
		{
			name: "item is out of bounds",
			n:    1e9,
			cmp: func(i int) int {
				if i >= 1e9 {
					return 0
				}
				return -1
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := xsort.Exist(tc.n, tc.cmp)

			if got != tc.expected {
				t.Errorf("expected %t; got %t", tc.expected, got)
			}
		})
	}
}

func TestExistInts(t *testing.T) {
	testCases := []struct {
		name     string
		a        []int
		x        int
		expected bool
	}{
		{
			name:     "empty",
			a:        nil,
			x:        0,
			expected: false,
		},
		{
			name:     "int is present",
			a:        []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			x:        7,
			expected: true,
		},
		{
			name:     "int is not present",
			a:        []int{0, 1, 2, 3, 4, 5, 6, 8, 9},
			x:        7,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := xsort.ExistInts(tc.a, tc.x)

			if got != tc.expected {
				t.Errorf("expected %t; got %t", tc.expected, got)
			}
		})
	}
}

func TestExistFloat64s(t *testing.T) {
	testCases := []struct {
		name     string
		a        []float64
		x        float64
		expected bool
	}{
		{
			name:     "empty",
			a:        nil,
			x:        0,
			expected: false,
		},
		{
			name:     "float64 is present",
			a:        []float64{0.1, 1.2, 2.3, 3.4, 4.5, 5.6, 6.7, 7.8, 8.9, 9.0},
			x:        7.8,
			expected: true,
		},
		{
			name:     "float64 is not present",
			a:        []float64{0.1, 1.2, 2.3, 3.4, 4.5, 5.6, 6.7, 7.8, 8.9, 9.0},
			x:        7.7,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := xsort.ExistFloat64s(tc.a, tc.x)

			if got != tc.expected {
				t.Errorf("expected %t; got %t", tc.expected, got)
			}
		})
	}
}

func TestExistStrings(t *testing.T) {
	testCases := []struct {
		name     string
		a        []string
		x        string
		expected bool
	}{
		{
			name:     "empty",
			a:        nil,
			x:        "f",
			expected: false,
		},
		{
			name:     "string is present",
			a:        []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			x:        "f",
			expected: true,
		},
		{
			name:     "string is not present",
			a:        []string{"a", "b", "c", "d", "e", "g", "h", "i", "j"},
			x:        "f",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := xsort.ExistStrings(tc.a, tc.x)

			if got != tc.expected {
				t.Errorf("expected %t; got %t", tc.expected, got)
			}
		})
	}
}
