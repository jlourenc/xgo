// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xsort_test

import (
	"fmt"

	"github.com/jlourenc/xgo/xsort"
)

// This example demonstrates checking a list sorted in ascending order.
func ExampleExist() {
	a := []int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55}
	x := 6

	exists := xsort.Exist(len(a), func(i int) int {
		if a[i] == x {
			return 0
		}
		if a[i] < x {
			return -1
		}
		return +1
	})
	if exists {
		fmt.Printf("element %d exists in %v\n", x, a)
	} else {
		fmt.Printf("element %d does not exist in %v\n", x, a)
	}
	// Output:
	// element 6 exists in [1 3 6 10 15 21 28 36 45 55]
}

// This example demonstrates checking a list sorted in descending order.
// The approach is the same as checking a list in ascending order,
// but with the condition inverted.
func ExampleExist_descendingOrder() {
	a := []int{55, 45, 36, 28, 21, 15, 10, 6, 3, 1}
	x := 6

	exists := xsort.Exist(len(a), func(i int) int {
		if a[i] == x {
			return 0
		}
		if a[i] > x {
			return -1
		}
		return +1
	})
	if exists {
		fmt.Printf("element %d exists in %v\n", x, a)
	} else {
		fmt.Printf("element %d does not exist in %v\n", x, a)
	}
	// Output:
	// element 6 exists in [55 45 36 28 21 15 10 6 3 1]
}

// This example demonstrates checking an int list sorted in ascending order.
func ExampleExistInts() {
	a := []int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55}
	x := 6

	exists := xsort.ExistInts(a, x)
	if exists {
		fmt.Printf("element %d exists in %v\n", x, a)
	} else {
		fmt.Printf("element %d does not exist in %v\n", x, a)
	}
	// Output:
	// element 6 exists in [1 3 6 10 15 21 28 36 45 55]
}

// This example demonstrates checking a float64 list sorted in ascending order.
func ExampleExistFloat64s() {
	a := []float64{1.5, 3.2, 6.9, 10.4, 15.7, 21.2, 28.9, 36.4, 45.8, 55.1}
	x := 6.9

	exists := xsort.ExistFloat64s(a, x)
	if exists {
		fmt.Printf("element %.2f exists in %v\n", x, a)
	} else {
		fmt.Printf("element %.2f does not exist in %v\n", x, a)
	}
	// Output:
	// element 6.90 exists in [1.5 3.2 6.9 10.4 15.7 21.2 28.9 36.4 45.8 55.1]
}

// This example demonstrates checking an int list sorted in ascending order.
func ExampleExistStrings() {
	a := []string{"aa", "efg", "i", "jjj", "lo", "mmn", "ok", "qts", "vw", "xyz"}
	x := "lo"

	exists := xsort.ExistStrings(a, x)
	if exists {
		fmt.Printf("element %q exists in %v\n", x, a)
	} else {
		fmt.Printf("element %q does not exist in %v\n", x, a)
	}
	// Output:
	// element "lo" exists in [aa efg i jjj lo mmn ok qts vw xyz]
}
