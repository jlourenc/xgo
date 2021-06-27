// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmath_test

import (
	"testing"

	. "github.com/jlourenc/xgo/xmath"
)

func TestDimInt(t *testing.T) {
	testCases := []struct {
		name     string
		x        int
		y        int
		expected int
	}{
		{
			name:     "zero values",
			x:        0,
			y:        0,
			expected: 0,
		},
		{
			name:     "x - y < 0 - positive values",
			x:        1,
			y:        2,
			expected: 0,
		},
		{
			name:     "x - y > 0 - positive values",
			x:        5,
			y:        2,
			expected: 3,
		},
		{
			name:     "x - y < 0 - negative values",
			x:        -5,
			y:        -2,
			expected: 0,
		},
		{
			name:     "x - y > 0 - negative values",
			x:        -1,
			y:        -2,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := DimInt(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestDimInt64(t *testing.T) {
	testCases := []struct {
		name     string
		x        int64
		y        int64
		expected int64
	}{
		{
			name:     "zero values",
			x:        0,
			y:        0,
			expected: 0,
		},
		{
			name:     "x - y < 0 - positive values",
			x:        1,
			y:        2,
			expected: 0,
		},
		{
			name:     "x - y > 0 - positive values",
			x:        5,
			y:        2,
			expected: 3,
		},
		{
			name:     "x - y < 0 - negative values",
			x:        -5,
			y:        -2,
			expected: 0,
		},
		{
			name:     "x - y > 0 - negative values",
			x:        -1,
			y:        -2,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := DimInt64(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestDimUint(t *testing.T) {
	testCases := []struct {
		name     string
		x        uint
		y        uint
		expected uint
	}{
		{
			name:     "zero values",
			x:        0,
			y:        0,
			expected: 0,
		},
		{
			name:     "x - y < 0",
			x:        1,
			y:        2,
			expected: 0,
		},
		{
			name:     "x - y > 0",
			x:        5,
			y:        2,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := DimUint(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestDimUint64(t *testing.T) {
	testCases := []struct {
		name     string
		x        uint64
		y        uint64
		expected uint64
	}{
		{
			name:     "zero values",
			x:        0,
			y:        0,
			expected: 0,
		},
		{
			name:     "x - y < 0",
			x:        1,
			y:        2,
			expected: 0,
		},
		{
			name:     "x - y > 0",
			x:        5,
			y:        2,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := DimUint64(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestMaxInt(t *testing.T) {
	testCases := []struct {
		name     string
		x        int
		y        int
		expected int
	}{
		{
			name:     "x == y",
			x:        2,
			y:        2,
			expected: 2,
		},
		{
			name:     "x < y",
			x:        1,
			y:        3,
			expected: 3,
		},
		{
			name:     "x > y",
			x:        3,
			y:        1,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaxInt(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestMaxInt64(t *testing.T) {
	testCases := []struct {
		name     string
		x        int64
		y        int64
		expected int64
	}{
		{
			name:     "x == y",
			x:        2,
			y:        2,
			expected: 2,
		},
		{
			name:     "x < y",
			x:        1,
			y:        3,
			expected: 3,
		},
		{
			name:     "x > y",
			x:        3,
			y:        1,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaxInt64(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestMaxUint(t *testing.T) {
	testCases := []struct {
		name     string
		x        uint
		y        uint
		expected uint
	}{
		{
			name:     "x == y",
			x:        2,
			y:        2,
			expected: 2,
		},
		{
			name:     "x < y",
			x:        1,
			y:        3,
			expected: 3,
		},
		{
			name:     "x > y",
			x:        3,
			y:        1,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaxUint(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestMaxUint64(t *testing.T) {
	testCases := []struct {
		name     string
		x        uint64
		y        uint64
		expected uint64
	}{
		{
			name:     "x == y",
			x:        2,
			y:        2,
			expected: 2,
		},
		{
			name:     "x < y",
			x:        1,
			y:        3,
			expected: 3,
		},
		{
			name:     "x > y",
			x:        3,
			y:        1,
			expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaxUint64(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestMinInt(t *testing.T) {
	testCases := []struct {
		name     string
		x        int
		y        int
		expected int
	}{
		{
			name:     "x == y",
			x:        2,
			y:        2,
			expected: 2,
		},
		{
			name:     "x < y",
			x:        1,
			y:        3,
			expected: 1,
		},
		{
			name:     "x > y",
			x:        3,
			y:        1,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MinInt(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestMinInt64(t *testing.T) {
	testCases := []struct {
		name     string
		x        int64
		y        int64
		expected int64
	}{
		{
			name:     "x == y",
			x:        2,
			y:        2,
			expected: 2,
		},
		{
			name:     "x < y",
			x:        1,
			y:        3,
			expected: 1,
		},
		{
			name:     "x > y",
			x:        3,
			y:        1,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MinInt64(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestMinUint(t *testing.T) {
	testCases := []struct {
		name     string
		x        uint
		y        uint
		expected uint
	}{
		{
			name:     "x == y",
			x:        2,
			y:        2,
			expected: 2,
		},
		{
			name:     "x < y",
			x:        1,
			y:        3,
			expected: 1,
		},
		{
			name:     "x > y",
			x:        3,
			y:        1,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MinUint(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestMinUint64(t *testing.T) {
	testCases := []struct {
		name     string
		x        uint64
		y        uint64
		expected uint64
	}{
		{
			name:     "x == y",
			x:        2,
			y:        2,
			expected: 2,
		},
		{
			name:     "x < y",
			x:        1,
			y:        3,
			expected: 1,
		},
		{
			name:     "x > y",
			x:        3,
			y:        1,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := MinUint64(tc.x, tc.y)

			if got != tc.expected {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}
