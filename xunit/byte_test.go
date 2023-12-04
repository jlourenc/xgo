// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xunit_test

import (
	"errors"
	"testing"

	"github.com/jlourenc/xgo/xunit"
)

func TestParseByte(t *testing.T) {
	testCases := []struct {
		input        string
		expectedByte xunit.Byte
		expectedErr  error
	}{
		{"", 0, errors.New("empty byte representation")},
		{"0.1.2KB", 0, errors.New("invalid byte representation: 0.1.2KB")},
		{"X", 0, errors.New("invalid byte representation: X")},
		{"1q", 0, errors.New("invalid byte representation: 1q")},
		{"9223372036854775808", 0, errors.New("invalid byte representation: 9223372036854775808")},
		{"-9223372036854775809", 0, errors.New("invalid byte representation: -9223372036854775809")},
		{"-9223372036854775808", -9223372036854775808, nil},
		{"-2PB", -2 * xunit.PB, nil},
		{"-2pb", -2 * xunit.PB, nil},
		{"-2Pb", -2 * xunit.PB, nil},
		{"-1TB", -xunit.TB, nil},
		{"-1TiB", -xunit.TiB, nil},
		{"-4096.75MiB", -4295753728, nil},
		{"-4096.5MiB", -4295491584, nil},
		{"-4096.25MiB", -4295229440, nil},
		{"-4096.000MiB", -4 * xunit.GiB, nil},
		{"-4096.0MiB", -4 * xunit.GiB, nil},
		{"-3GiB", -3 * xunit.GiB, nil},
		{"-3GB", -3 * xunit.GB, nil},
		{"-1.5GiB", -1610612736, nil},
		{"-1.5GB", -1500000000, nil},
		{"-5MiB", -5 * xunit.MiB, nil},
		{"-5MB", -5 * xunit.MB, nil},
		{"-1.5Mib", -1572864, nil},
		{"-10KiB", -10 * xunit.KiB, nil},
		{"-10KB", -10 * xunit.KB, nil},
		{"-1B", -xunit.B, nil},
		{"-0", 0, nil},
		{"0", 0, nil},
		{"1B", xunit.B, nil},
		{"10KB", 10 * xunit.KB, nil},
		{"10KiB", 10 * xunit.KiB, nil},
		{"1.5Mib", 1572864, nil},
		{"5MB", 5 * xunit.MB, nil},
		{"5MiB", 5 * xunit.MiB, nil},
		{"1.5GB", 1500000000, nil},
		{"1.5GiB", 1610612736, nil},
		{"3GB", 3 * xunit.GB, nil},
		{"3GiB", 3 * xunit.GiB, nil},
		{"4096.0MiB", 4 * xunit.GiB, nil},
		{"4096.000MiB", 4 * xunit.GiB, nil},
		{"4096.25MiB", 4295229440, nil},
		{"4096.5MiB", 4295491584, nil},
		{"4096.75MiB", 4295753728, nil},
		{"1TiB", xunit.TiB, nil},
		{"1TB", xunit.TB, nil},
		{"2Pb", 2 * xunit.PB, nil},
		{"2pb", 2 * xunit.PB, nil},
		{"2PB", 2 * xunit.PB, nil},
		{"9223372036854775807", 9223372036854775807, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			qty, err := xunit.ParseByte(tc.input)

			if tc.expectedByte != qty {
				t.Errorf("expected %s; got %s", tc.expectedByte, qty)
			}

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) ||
				(tc.expectedErr != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("expected error %s; got %s", tc.expectedErr, err)
			}
		})
	}
}

func TestByte_Get(t *testing.T) {
	b := 2*xunit.MiB + 512*xunit.KiB

	got := b.Get()

	if got != b {
		t.Errorf("expected %s; got %s", b, got)
	}
}

func TestByte_Set(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedByte xunit.Byte
		expectedErr  error
	}{
		{
			name:        "empty byte representation",
			input:       "",
			expectedErr: errors.New("empty byte representation"),
		},
		{
			name:        "invalid byte representation",
			input:       "2X",
			expectedErr: errors.New("invalid byte representation: 2X"),
		},
		{
			name:         "valid byte representation",
			input:        "2MiB",
			expectedByte: 2 * xunit.MiB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b xunit.Byte

			err := b.Set(tc.input)

			if tc.expectedByte != b {
				t.Errorf("expected %s; got %s", tc.expectedByte, b)
			}

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) ||
				(tc.expectedErr != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("expected error %s; got %s", tc.expectedErr, err)
			}
		})
	}
}

func TestByte_MarshalText_String(t *testing.T) {
	testCases := []struct {
		name     string
		input    xunit.Byte
		expected string
	}{
		{"-9223372036854775808", -9223372036854775808, "-8EiB"},
		{"-2PB", -2 * xunit.PB, "-2PB"},
		{"-10TiB", -10 * xunit.TiB, "-10TiB"},
		{"-1TB", -xunit.TB, "-1TB"},
		{"-4096.75MiB", -4295753728, "-4096.75MiB"},
		{"-4096.5MiB", -4295491584, "-4096.5MiB"},
		{"-4096.25MiB", -4295229440, "-4096.25MiB"},
		{"-4096.0MiB", -4 * xunit.GiB, "-4GiB"},
		{"-4096.000MiB", -4 * xunit.GiB, "-4GiB"},
		{"-1.5GiB", -1610612736, "-1.5GiB"},
		{"-1.5GB", -1500000000, "-1.5GB"},
		{"-1GiB", -xunit.GiB, "-1GiB"},
		{"-1GB", -xunit.GB, "-1GB"},
		{"-1.5Mib", -1572864, "-1.5MiB"},
		{"-1MiB", -xunit.MiB, "-1MiB"},
		{"-1MB", -xunit.MB, "-1MB"},
		{"-1KiB", -xunit.KiB, "-1KiB"},
		{"-1KB", -xunit.KB, "-1KB"},
		{"-1B", -xunit.B, "-1B"},
		{"-0", 0, "0B"},
		{"0", 0, "0B"},
		{"1B", xunit.B, "1B"},
		{"1KB", xunit.KB, "1KB"},
		{"1KiB", xunit.KiB, "1KiB"},
		{"1MB", xunit.MB, "1MB"},
		{"1MiB", xunit.MiB, "1MiB"},
		{"512MiB", 512 * xunit.MiB, "512MiB"},
		{"768MiB", 768 * xunit.MiB, "768MiB"},
		{"1.5Mib", 1572864, "1.5MiB"},
		{"1GB", xunit.GB, "1GB"},
		{"1GiB", xunit.GiB, "1GiB"},
		{"1.5GB", 1500000000, "1.5GB"},
		{"1.5GiB", 1610612736, "1.5GiB"},
		{"4096.000MiB", 4 * xunit.GiB, "4GiB"},
		{"4096.0MiB", 4 * xunit.GiB, "4GiB"},
		{"4096.25MiB", 4295229440, "4096.25MiB"},
		{"4096.5MiB", 4295491584, "4096.5MiB"},
		{"4096.75MiB", 4295753728, "4096.75MiB"},
		{"1TB", xunit.TB, "1TB"},
		{"10TiB", 10 * xunit.TiB, "10TiB"},
		{"2PB", 2 * xunit.PB, "2PB"},
		{"9223372036854775807", 9223372036854775807, "8EiB"},
	}

	for _, tc := range testCases {
		t.Run(tc.name+"_string", func(t *testing.T) {
			got := tc.input.String()

			if tc.expected != got {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
		t.Run(tc.name+"_marshal_text", func(t *testing.T) {
			got, err := tc.input.MarshalText()

			if tc.expected != string(got) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}

			if err != nil {
				t.Errorf("no error expected; got %s", err)
			}
		})
	}
}

func TestByte_Type(t *testing.T) {
	var b xunit.Byte
	expected := "xunit_byte"

	got := b.Type()

	if got != expected {
		t.Errorf("expected %s; got %s", expected, got)
	}
}

func TestByte_UnmarshalText(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expectedByte xunit.Byte
		expectedErr  error
	}{
		{
			name:        "empty byte representation",
			input:       "",
			expectedErr: errors.New("empty byte representation"),
		},
		{
			name:        "invalid byte representation",
			input:       "2X",
			expectedErr: errors.New("invalid byte representation: 2X"),
		},
		{
			name:         "valid byte representation",
			input:        "2MiB",
			expectedByte: 2 * xunit.MiB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b xunit.Byte

			err := b.UnmarshalText([]byte(tc.input))

			if tc.expectedByte != b {
				t.Errorf("expected %s; got %s", tc.expectedByte, b)
			}

			if (tc.expectedErr == nil && err != nil) || (tc.expectedErr != nil && err == nil) ||
				(tc.expectedErr != nil && tc.expectedErr.Error() != err.Error()) {
				t.Errorf("expected error %s; got %s", tc.expectedErr, err)
			}
		})
	}
}

func TestByte_B(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{input: "-2.5KB", expected: -2500},
		{input: "-1.5KiB", expected: -1536},
		{input: "-1B", expected: -1},
		{input: "1B", expected: 1},
		{input: "1.5KiB", expected: 1536},
		{input: "2.5KB", expected: 2500},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.B()

			if tc.expected != got {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestByte_KB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5KB", expected: -2.5},
		{input: "-2KB", expected: -2},
		{input: "-1500B", expected: -1.5},
		{input: "-1000B", expected: -1},
		{input: "1000B", expected: 1},
		{input: "1500B", expected: 1.5},
		{input: "2KB", expected: 2},
		{input: "2.5KB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.KB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_KiB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5KiB", expected: -2.5},
		{input: "-2KiB", expected: -2},
		{input: "-1536B", expected: -1.5},
		{input: "-1024B", expected: -1},
		{input: "1024B", expected: 1},
		{input: "1536B", expected: 1.5},
		{input: "2KiB", expected: 2},
		{input: "2.5KiB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.KiB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_MB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5MB", expected: -2.5},
		{input: "-2MB", expected: -2},
		{input: "-1500KB", expected: -1.5},
		{input: "-1000KB", expected: -1},
		{input: "1000KB", expected: 1},
		{input: "1500KB", expected: 1.5},
		{input: "2MB", expected: 2},
		{input: "2.5MB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.MB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_MiB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5MiB", expected: -2.5},
		{input: "-2MiB", expected: -2},
		{input: "-1536KiB", expected: -1.5},
		{input: "-1024KiB", expected: -1},
		{input: "1024KiB", expected: 1},
		{input: "1536KiB", expected: 1.5},
		{input: "2MiB", expected: 2},
		{input: "2.5MiB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.MiB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_GB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5GB", expected: -2.5},
		{input: "-2GB", expected: -2},
		{input: "-1500MB", expected: -1.5},
		{input: "-1000MB", expected: -1},
		{input: "1000MB", expected: 1},
		{input: "1500MB", expected: 1.5},
		{input: "2GB", expected: 2},
		{input: "2.5GB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.GB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_GiB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5GiB", expected: -2.5},
		{input: "-2GiB", expected: -2},
		{input: "-1536MiB", expected: -1.5},
		{input: "-1024MiB", expected: -1},
		{input: "1024MiB", expected: 1},
		{input: "1536MiB", expected: 1.5},
		{input: "2GiB", expected: 2},
		{input: "2.5GiB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.GiB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_TB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5TB", expected: -2.5},
		{input: "-2TB", expected: -2},
		{input: "-1500GB", expected: -1.5},
		{input: "-1000GB", expected: -1},
		{input: "1000GB", expected: 1},
		{input: "1500GB", expected: 1.5},
		{input: "2TB", expected: 2},
		{input: "2.5TB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.TB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_TiB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5TiB", expected: -2.5},
		{input: "-2TiB", expected: -2},
		{input: "-1536GiB", expected: -1.5},
		{input: "-1024GiB", expected: -1},
		{input: "1024GiB", expected: 1},
		{input: "1536GiB", expected: 1.5},
		{input: "2TiB", expected: 2},
		{input: "2.5TiB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.TiB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_PB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5PB", expected: -2.5},
		{input: "-2PB", expected: -2},
		{input: "-1500TB", expected: -1.5},
		{input: "-1000TB", expected: -1},
		{input: "1000TB", expected: 1},
		{input: "1500TB", expected: 1.5},
		{input: "2PB", expected: 2},
		{input: "2.5PB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.PB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_PiB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5PiB", expected: -2.5},
		{input: "-2PiB", expected: -2},
		{input: "-1536TiB", expected: -1.5},
		{input: "-1024TiB", expected: -1},
		{input: "1024TiB", expected: 1},
		{input: "1536TiB", expected: 1.5},
		{input: "2PiB", expected: 2},
		{input: "2.5PiB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.PiB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_EB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5EB", expected: -2.5},
		{input: "-2EB", expected: -2},
		{input: "-1500PB", expected: -1.5},
		{input: "-1000PB", expected: -1},
		{input: "1000PB", expected: 1},
		{input: "1500PB", expected: 1.5},
		{input: "2EB", expected: 2},
		{input: "2.5EB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.EB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}

func TestByte_EiB(t *testing.T) {
	testCases := []struct {
		input    string
		expected float64
	}{
		{input: "-2.5EiB", expected: -2.5},
		{input: "-2EiB", expected: -2},
		{input: "-1536PiB", expected: -1.5},
		{input: "-1024PiB", expected: -1},
		{input: "1024PiB", expected: 1},
		{input: "1536PiB", expected: 1.5},
		{input: "2EiB", expected: 2},
		{input: "2.5EiB", expected: 2.5},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			b, _ := xunit.ParseByte(tc.input)
			got := b.EiB()

			if tc.expected != got {
				t.Errorf("expected %f; got %f", tc.expected, got)
			}
		})
	}
}
