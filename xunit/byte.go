// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xunit extends the Go standard library by providing
// additional primitives and structures for mainpulating certain units.
package xunit

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func init() {
	for k, v := range byteUnits {
		byteUnitsString[strings.ToLower(v)] = k
	}
}

// Enumeration of byte units.
const (
	B Byte = 1

	// Decimal.
	KB = B * 1000  // 10^3 bytes
	MB = KB * 1000 // 10^6 bytes
	GB = MB * 1000 // 10^9 bytes
	TB = GB * 1000 // 10^12 bytes
	PB = TB * 1000 // 10^15 bytes
	EB = PB * 1000 // 10^18 bytes

	// Binary.
	KiB = B << 10   // 2^10 bytes
	MiB = KiB << 10 // 2^20 bytes
	GiB = MiB << 10 // 2^30 bytes
	TiB = GiB << 10 // 2^40 bytes
	PiB = TiB << 10 // 2^50 bytes
	EiB = PiB << 10 // 2^60 bytes
)

const (
	errByteEmptyMsg   = "empty byte representation"
	errByteInvalidMsg = "invalid byte representation: "
)

// Byte is a count of bytes.
type Byte int64

var (
	byteUnits = map[Byte]string{
		B:   "B",
		KB:  "KB",
		KiB: "KiB",
		MB:  "MB",
		MiB: "MiB",
		GB:  "GB",
		GiB: "GiB",
		TB:  "TB",
		TiB: "TiB",
		PB:  "PB",
		PiB: "PiB",
		EB:  "EB",
		EiB: "EiB",
	}
	byteUnitsString = map[string]Byte{
		"": B,
	}
	bytesUnitsDescOrder = []Byte{EiB, EB, PiB, PB, TiB, TB, GiB, GB, MiB, MB, KiB, KB}
)

// Parse parses a byte string which is a number followed by a byte unit suffix (e.g. '1024MB' or '1GiB').
// The following units are available:
//
//	B:   Byte
//	KB:  Kilobyte
//	KiB: Kibibyte
//	MB:  Megabyte
//	MiB: Mebibyte
//	GB:  Gigabyte
//	GiB: Gibibyte
//	TB:  Terabyte
//	TiB: Tebibyte
//	PB:  Petabyte
//	PiB: Pebibbyte
//	EB:  Exabyte
//	EiB: Exbibyte
func ParseByte(s string) (Byte, error) {
	s = strings.TrimSpace(s)

	if s == "" {
		return 0, errors.New(errByteEmptyMsg)
	}

	isFloat := false
	i := 0

strLoop:
	for _, r := range s {
		switch {
		case r == '.':
			isFloat = true
		case !unicode.IsDigit(r) && r != '-':
			break strLoop
		}
		i++
	}

	unit, ok := byteUnitsString[strings.ToLower(s[i:])]
	if !ok {
		return 0, errors.New(errByteInvalidMsg + s)
	}

	if !isFloat { // no fractional floating-point numbers
		qty, err := strconv.ParseInt(s[:i], 10, 64)
		if err != nil {
			return 0, errors.New(errByteInvalidMsg + s)
		}
		return Byte(qty) * unit, nil
	}

	qty, err := strconv.ParseFloat(s[:i], 64)
	if err != nil {
		return 0, errors.New(errByteInvalidMsg + s)
	}

	whole, frac := math.Modf(qty)
	return Byte((whole * float64(unit)) + (frac * float64(unit))), nil
}

// Get returns the Byte value.
// It makes Byte implement the flag package Getter interface.
func (b Byte) Get() any { return b }

// MarshalText implements the encoding.TextMarshaler interface.
// The encoding is the same as returned by String.
func (b Byte) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

// Set parses the string in input and assign it to b if valid, otherwise an error is returned.
// It makes Byte implement the flag package Value interface.
func (b *Byte) Set(s string) error {
	bs, err := ParseByte(s)
	if err != nil {
		return err
	}
	*b = bs
	return nil
}

// String returns a string representation of Byte with the most suitable unit.
func (b Byte) String() string {
	if b == 0 {
		return "0B"
	}

	for _, unit := range bytesUnitsDescOrder {
		qty := b.toUnit(unit)

		if math.Abs(qty) < 1 {
			continue
		}

		if checkDecimalPlaces(0, qty) {
			return strconv.FormatInt(int64(qty), 10) + byteUnits[unit]
		}

		if checkDecimalPlaces(2, qty) {
			return strconv.FormatFloat(qty, 'g', 53, 64) + byteUnits[unit]
		}
	}

	return strconv.FormatInt(int64(b), 10) + byteUnits[B]
}

// Type returns a string representation of Byte type.
// It makes Byte implement the pflag Value interface.
func (Byte) Type() string { return "xunit_byte" }

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The text is expected in a form accepted by ParseByte.
func (b *Byte) UnmarshalText(text []byte) error {
	return b.Set(string(text))
}

// B returns the value in bytes.
func (b Byte) B() int64 {
	return int64(b)
}

// KB returns the value in kilobytes (10^3).
func (b Byte) KB() float64 {
	return b.toUnit(KB)
}

// MB returns the value in megabytes (10^6).
func (b Byte) MB() float64 {
	return b.toUnit(MB)
}

// GB returns the value in gigabytes (10^9).
func (b Byte) GB() float64 {
	return b.toUnit(GB)
}

// TB returns the value in terabytes (10^12).
func (b Byte) TB() float64 {
	return b.toUnit(TB)
}

// PB returns the value in petabytes (10^15).
func (b Byte) PB() float64 {
	return b.toUnit(PB)
}

// EB returns the value in exabytes (10^18).
func (b Byte) EB() float64 {
	return b.toUnit(EB)
}

// KiB returns the value in kibibytes (2^10).
func (b Byte) KiB() float64 {
	return b.toUnit(KiB)
}

// MiB returns the value in mebibytes (2^20).
func (b Byte) MiB() float64 {
	return b.toUnit(MiB)
}

// GiB returns the value in gibibytes (2^30).
func (b Byte) GiB() float64 {
	return b.toUnit(GiB)
}

// TiB returns the value in tebibytes (2^40).
func (b Byte) TiB() float64 {
	return b.toUnit(TiB)
}

// PiB returns the value in pebibytes (2^50).
func (b Byte) PiB() float64 {
	return b.toUnit(PiB)
}

// EiB returns the value in exbibytes (2^60).
func (b Byte) EiB() float64 {
	return b.toUnit(EiB)
}

func (b Byte) toUnit(unit Byte) float64 {
	whole := b / unit
	remainder := b - (whole * unit)
	return float64(whole) + float64(remainder)/float64(unit)
}

func checkDecimalPlaces(i int, value float64) bool {
	value *= math.Pow(10.0, float64(i))
	extra := value - float64(int64(value))
	return extra == 0
}
