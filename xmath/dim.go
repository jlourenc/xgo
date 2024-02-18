// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xmath extends the Go standard library package math
// by providing additional mathematical functions.
package xmath

// DimInt returns the maximum of x-y or 0
// with x, y and the return value of type int.
func DimInt(x, y int) int {
	v := x - y
	if v <= 0 {
		return 0
	}
	return v
}

// DimInt64 returns the maximum of x-y or 0
// with x, y and the return value of type int64.
func DimInt64(x, y int64) int64 {
	v := x - y
	if v <= 0 {
		return 0
	}
	return v
}

// DimUint returns the maximum of x-y or 0
// with x, y and the return value of type uint.
func DimUint(x, y uint) uint {
	if x > y {
		return x - y
	}
	return 0
}

// DimUint64 returns the maximum of x-y or 0
// with x, y and the return value of type uint64.
func DimUint64(x, y uint64) uint64 {
	if x > y {
		return x - y
	}
	return 0
}

// MaxInt returns the larger of x or y
// with x, y and the return value of type int.
//
// Deprecated: From Go 1.21, use the built-in max function instead.
func MaxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// MaxInt64 returns the larger of x or y
// with x, y and the return value of type int64.
//
// Deprecated: From Go 1.21, use the built-in max function instead.
func MaxInt64(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// MaxUint returns the larger of x or y
// with x, y and the return value of type uint.
//
// Deprecated: From Go 1.21, use the built-in max function instead.
func MaxUint(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

// MaxUint64 returns the larger of x or y
// with x, y and the return value of type uint64.
//
// Deprecated: From Go 1.21, use the built-in max function instead.
func MaxUint64(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}

// MinInt returns the smaller of x or y
// with x, y and the return value of type int.
//
// Deprecated: From Go 1.21, use the built-in min function instead.
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// MinInt64 returns the smaller of x or y
// with x, y and the return value of type int64.
//
// Deprecated: From Go 1.21, use the built-in min function instead.
func MinInt64(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// MinUint returns the smaller of x or y
// with x, y and the return value of type uint.
//
// Deprecated: From Go 1.21, use the built-in min function instead.
func MinUint(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}

// MinUint64 returns the smaller of x or y
// with x, y and the return value of type uint64.
//
// Deprecated: From Go 1.21, use the built-in min function instead.
func MinUint64(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
