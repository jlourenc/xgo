// Copyright 2022 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xmath_test

import (
	"fmt"

	"github.com/jlourenc/xgo/xmath"
)

func ExampleDimInt() {
	fmt.Printf("%d\n", xmath.DimInt(4, -2))
	fmt.Printf("%d\n", xmath.DimInt(-4, 2))
	// Output:
	// 6
	// 0
}

func ExampleDimInt64() {
	fmt.Printf("%d\n", xmath.DimInt64(4, -2))
	fmt.Printf("%d\n", xmath.DimInt64(-4, 2))
	// Output:
	// 6
	// 0
}

func ExampleDimUint() {
	fmt.Printf("%d\n", xmath.DimUint(4, 2))
	fmt.Printf("%d\n", xmath.DimUint(4, 6))
	// Output:
	// 2
	// 0
}

func ExampleDimUint64() {
	fmt.Printf("%d\n", xmath.DimUint64(4, 2))
	fmt.Printf("%d\n", xmath.DimUint64(4, 6))
	// Output:
	// 2
	// 0
}

func ExampleMaxInt() {
	fmt.Printf("%d\n", xmath.MaxInt(4, 2))
	fmt.Printf("%d\n", xmath.MaxInt(4, -2))
	fmt.Printf("%d\n", xmath.MaxInt(-4, -2))
	// Output:
	// 4
	// 4
	// -2
}

func ExampleMaxInt64() {
	fmt.Printf("%d\n", xmath.MaxInt64(4, 2))
	fmt.Printf("%d\n", xmath.MaxInt64(4, -2))
	fmt.Printf("%d\n", xmath.MaxInt64(-4, -2))
	// Output:
	// 4
	// 4
	// -2
}

func ExampleMaxUint() {
	fmt.Printf("%d\n", xmath.MaxUint(4, 2))
	// Output:
	// 4
}

func ExampleMaxUint64() {
	fmt.Printf("%d\n", xmath.MaxUint64(4, 2))
	// Output:
	// 4
}

func ExampleMinInt() {
	fmt.Printf("%d\n", xmath.MinInt(4, 2))
	fmt.Printf("%d\n", xmath.MinInt(4, -2))
	fmt.Printf("%d\n", xmath.MinInt(-4, -2))
	// Output:
	// 2
	// -2
	// -4
}

func ExampleMinInt64() {
	fmt.Printf("%d\n", xmath.MinInt64(4, 2))
	fmt.Printf("%d\n", xmath.MinInt64(4, -2))
	fmt.Printf("%d\n", xmath.MinInt64(-4, -2))
	// Output:
	// 2
	// -2
	// -4
}

func ExampleMinUint() {
	fmt.Printf("%d\n", xmath.MinUint(4, 2))
	// Output:
	// 2
}

func ExampleMinUint64() {
	fmt.Printf("%d\n", xmath.MinUint64(4, 2))
	// Output:
	// 2
}
