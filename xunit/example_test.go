// Copyright 2022 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xunit_test

import (
	"fmt"

	"github.com/jlourenc/xgo/xunit"
)

func ExampleParseByte() {
	b, err := xunit.ParseByte("2048MiB")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", b)
	// Output: 2GiB
}

func ExampleByte_MarshalText() {
	b := xunit.TiB + 512*xunit.GiB
	bytes, err := b.MarshalText()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", bytes)
	// Output: 1.5TiB
}

func ExampleByte_Set() {
	var b xunit.Byte
	if err := b.Set("512MiB"); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", b)
	// Output: 512MiB
}

func ExampleByte_UnmarshalText() {
	var b xunit.Byte
	if err := b.UnmarshalText([]byte("512PiB")); err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", b)
	// Output: 512PiB
}

func ExampleByte_B() {
	b := 1 * xunit.EB
	fmt.Printf("%d\n", b.B())
	// Output: 1000000000000000000
}

func ExampleByte_KB() {
	b := 1 * xunit.EB
	fmt.Printf("%.0f\n", b.KB())
	// Output: 1000000000000000
}

func ExampleByte_MB() {
	b := 1 * xunit.EB
	fmt.Printf("%.0f\n", b.MB())
	// Output: 1000000000000
}

func ExampleByte_GB() {
	b := 1 * xunit.EB
	fmt.Printf("%.0f\n", b.GB())
	// Output: 1000000000
}

func ExampleByte_TB() {
	b := 1 * xunit.EB
	fmt.Printf("%.0f\n", b.TB())
	// Output: 1000000
}

func ExampleByte_PB() {
	b := 1 * xunit.EB
	fmt.Printf("%.0f\n", b.PB())
	// Output: 1000
}

func ExampleByte_EB() {
	b := 1 * xunit.EB
	fmt.Printf("%.0f\n", b.EB())
	// Output: 1
}

func ExampleByte_KiB() {
	b := 1 * xunit.EiB
	fmt.Printf("%.0f\n", b.KiB())
	// Output: 1125899906842624
}

func ExampleByte_MiB() {
	b := 1 * xunit.EiB
	fmt.Printf("%.0f\n", b.MiB())
	// Output: 1099511627776
}

func ExampleByte_GiB() {
	b := 1 * xunit.EiB
	fmt.Printf("%.0f\n", b.GiB())
	// Output: 1073741824
}

func ExampleByte_TiB() {
	b := 1 * xunit.EiB
	fmt.Printf("%.0f\n", b.TiB())
	// Output: 1048576
}

func ExampleByte_PiB() {
	b := 1 * xunit.EiB
	fmt.Printf("%.0f\n", b.PiB())
	// Output: 1024
}

func ExampleByte_EiB() {
	b := 1 * xunit.EiB
	fmt.Printf("%.0f\n", b.EiB())
	// Output: 1
}
