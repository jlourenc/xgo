// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsort

const (
	minSizeForBinarySearch = 6 // minimum size by which binary search becomes more performant than linear search
)

// Exist uses either linear search (for slices with less than 6 elements) or binary search (otherwise)
// to check if there is an index i for which cmp(i) returns 0, meaning a value x exists in
// a sorted, indexable data structure such as an array or slice. In this case, the argument cmp,
// typically a closure, captures the value to be searched for, and how the data structure is indexed
// and ordered.
//
// For a given slice data sorted in ascending order:
// * Compare function cmp should return -1 if the value at index i is lower than the searched value,
// * Compare function cmp should return +1 if the value at index i is greater than the searched value.
//
// On the contrary, for a given slice data sorted in descending order:
// * Compare function cmp should return -1 if the value at index i is greater than the searched value,
// * Compare function cmp should return +1 if the value at index i is lower than the searched value.
//
// Exist calls cmp(i) only for i in the range [0, n).
//
// See ExistInts for an example of usage.
//
func Exist(n int, cmp func(int) int) bool {
	if n < minSizeForBinarySearch { // linear search
		for i := 0; i < n; i++ {
			if cmp(i) == 0 {
				return true
			}
		}
		return false
	}

	// binary search
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h

		switch cmp(h) {
		case 0:
			return true
		case -1:
			i = h + 1
		case +1:
			j = h
		}
	}
	return false
}

// Convenience wrappers for common cases.

// ExistInts returns whether or not x exists in a sorted slice of ints.
// The slice must be sorted in ascending order.
//
func ExistInts(a []int, x int) bool {
	return Exist(len(a), func(i int) int {
		if a[i] == x {
			return 0
		}
		if a[i] < x {
			return -1
		}
		return +1
	})
}

// ExistFloat64s returns whether or not x exists in a sorted slice of float64s.
// The slice must be sorted in ascending order.
//
func ExistFloat64s(a []float64, x float64) bool {
	return Exist(len(a), func(i int) int {
		if a[i] == x {
			return 0
		}
		if a[i] < x {
			return -1
		}
		return +1
	})
}

// ExistStrings returns whether or not x exists in a sorted slice of strings.
// The slice must be sorted in ascending order.
//
func ExistStrings(a []string, x string) bool {
	return Exist(len(a), func(i int) int {
		if a[i] == x {
			return 0
		}
		if a[i] < x {
			return -1
		}
		return +1
	})
}
