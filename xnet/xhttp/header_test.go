// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xhttp_test

import (
	"net/http"
	"testing"
	"time"

	. "github.com/jlourenc/xgo/xnet/xhttp"
)

func TestHeaderExist(t *testing.T) {
	testCases := []struct {
		name     string
		headers  http.Header
		key      string
		expected bool
	}{
		{
			name: "canonicalized key exists",
			headers: http.Header{
				"Header-Key": []string{},
			},
			key:      "Header-Key",
			expected: true,
		},
		{
			name: "non-canonicalized key exists",
			headers: http.Header{
				"Header-Key": []string{},
			},
			key:      "header-key",
			expected: true,
		},
		{
			name: "key does not exist",
			headers: http.Header{
				"Header-Key": []string{},
			},
			key:      "unknown",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := HeaderExist(tc.headers, tc.key)

			if tc.expected != got {
				t.Fatalf("expected %v; got %v", tc.expected, got)
			}
		})
	}
}

func TestHeaderKeyValues(t *testing.T) {
	testCases := []struct {
		name     string
		headers  http.Header
		key      string
		expected map[string]string
	}{
		{
			name:     "nil headers",
			headers:  nil,
			key:      "Header-Key",
			expected: nil,
		}, {
			name: "nil header values",
			headers: http.Header{
				"Header-Key": nil,
			},
			key:      "Header-Key",
			expected: nil,
		}, {
			name: "empty header values",
			headers: http.Header{
				"Header-Key": {},
			},
			key:      "Header-Key",
			expected: map[string]string{},
		}, {
			name: "single header value with single key",
			headers: http.Header{
				"Header-Key": {"key1"},
			},
			key: "Header-Key",
			expected: map[string]string{
				"key1": "",
			},
		}, {
			name: "single header value with single key/value pair",
			headers: http.Header{
				"Header-Key": {"key1=val1"},
			},
			key: "Header-Key",
			expected: map[string]string{
				"key1": "val1",
			},
		}, {
			name: "multiple header values with multiple keys or key/value pairs",
			headers: http.Header{
				"Header-Key": {"key1=val1", "key2", "key3=val3, key4, key5=val5"},
			},
			key: "Header-Key",
			expected: map[string]string{
				"key1": "val1",
				"key2": "",
				"key3": "val3",
				"key4": "",
				"key5": "val5",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := HeaderKeyValues(tc.headers, tc.key)

			if len(tc.expected) != len(got) {
				t.Fatalf("expected %v; got %v", tc.expected, got)
			}

			for key, value := range tc.expected {
				if v, ok := got[key]; !ok || value != v {
					t.Fatalf("expected %v; got %v", tc.expected, got)
				}
			}
		})
	}
}

func TestHeaderValues(t *testing.T) {
	testCases := []struct {
		name     string
		headers  http.Header
		key      string
		expected []string
	}{
		{
			name:     "nil header",
			headers:  nil,
			key:      "Header-Key",
			expected: nil,
		}, {
			name: "nil header values",
			headers: http.Header{
				"Header-Key": nil,
			},
			key:      "Header-Key",
			expected: nil,
		}, {
			name: "empty header values",
			headers: http.Header{
				"Header-Key": {},
			},
			key:      "Header-Key",
			expected: []string{},
		}, {
			name: "single header value with single field",
			headers: http.Header{
				"Header-Key": {"val1"},
			},
			key:      "Header-Key",
			expected: []string{"val1"},
		}, {
			name: "single header value with multiple fields",
			headers: http.Header{
				"Header-Key": {"val1, val2"},
			},
			key:      "Header-Key",
			expected: []string{"val1", "val2"},
		}, {
			name: "multiple values with multiple fields",
			headers: http.Header{
				"Header-Key": {"val1, val2", "val3", "val4, val5"},
			},
			key:      "Header-Key",
			expected: []string{"val1", "val2", "val3", "val4", "val5"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := HeaderValues(tc.headers, tc.key)

			if len(tc.expected) != len(got) {
				t.Fatalf("expected %v; got %v", tc.expected, got)
			}

			for i, value := range tc.expected {
				if value != got[i] {
					t.Fatalf("expected %v; got %v", tc.expected, got)
				}
			}
		})
	}
}

func TestParseHeaderDate(t *testing.T) {
	testCases := []struct {
		name         string
		headers      http.Header
		expectedTime time.Time
		expectedErr  bool
	}{
		{
			name:         "undefined",
			headers:      nil,
			expectedTime: time.Time{},
			expectedErr:  true,
		},
		{
			name: "empty",
			headers: http.Header{
				HeaderDate: {""},
			},
			expectedTime: time.Time{},
			expectedErr:  true,
		},
		{
			name: "invalid date",
			headers: http.Header{
				HeaderDate: {"invalid"},
			},
			expectedTime: time.Time{},
			expectedErr:  true,
		},
		{
			name: "invalid format",
			headers: http.Header{
				HeaderDate: {"1994-11-06T08:49:37Z00:00"},
			},
			expectedTime: time.Time{},
			expectedErr:  true,
		},
		{
			name: "RFC1123 format",
			headers: http.Header{
				HeaderDate: {"Sun, 10 Jul 2016 21:12:00.499 GMT"},
			},
			expectedTime: time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.UTC),
			expectedErr:  false,
		},
		{
			name: "RFC850 format",
			headers: http.Header{
				HeaderDate: {"Sunday, 10-Jul-16 21:12:00 GMT"},
			},
			expectedTime: time.Date(2016, time.July, 10, 21, 12, 0, 0, time.UTC),
			expectedErr:  false,
		},
		{
			name: "ANSIC format",
			headers: http.Header{
				HeaderDate: {"Sun Jul 10 21:12:00 2016"},
			},
			expectedTime: time.Date(2016, time.July, 10, 21, 12, 0, 0, time.UTC),
			expectedErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := ParseHeaderDate(tc.headers)

			if tc.expectedErr && err == nil {
				t.Error("error expected, got nil")
			} else if !tc.expectedErr && err != nil {
				t.Errorf("no error expected, got %s", err)
			}

			if !tc.expectedTime.Equal(d) {
				t.Errorf("expected %s, got %s", tc.expectedTime, d)
			}
		})
	}
}

func TestReplaceHeader(t *testing.T) {
	testCases := []struct {
		name     string
		headers  http.Header
		prefix   string
		key      string
		values   []string
		expected http.Header
	}{
		{
			name:     "nil headers",
			headers:  nil,
			prefix:   "Prefix",
			key:      "Header-Key",
			values:   []string{"val1", "val2", "val3"},
			expected: nil,
		},
		{
			name:    "empty headers - no values",
			headers: http.Header{},
			prefix:  "Prefix",
			key:     "Header-Key",
			values:  nil,
			expected: http.Header{
				"Header-Key": []string{},
			},
		},
		{
			name:    "empty headers - multiples values",
			headers: http.Header{},
			prefix:  "Prefix",
			key:     "Header-Key",
			values:  []string{"val1, val2", "val3", "val4"},
			expected: http.Header{
				"Header-Key": []string{"val1, val2", "val3", "val4"},
			},
		},
		{
			name: "multiple headers - multiples values",
			headers: http.Header{
				"Header-Key":      []string{"val5, val6", "val7", "val8"},
				"Header-Key-Copy": []string{"val5, val6", "val7", "val8"},
			},
			prefix: "Prefix",
			key:    "Header-Key",
			values: []string{"val1, val2", "val3", "val4"},
			expected: http.Header{
				"Header-Key":        []string{"val1, val2", "val3", "val4"},
				"Prefix-Header-Key": []string{"val5, val6", "val7", "val8"},
				"Header-Key-Copy":   []string{"val5, val6", "val7", "val8"},
			},
		},
		{
			name: "multiple prefixed headers - multiples values",
			headers: http.Header{
				"Prefix-Header-Key":   []string{"val50, val60", "val70", "val80"},
				"Prefix-1-Header-Key": []string{"val51, val61", "val71", "val81"},
				"Prefix-2-Header-Key": []string{"val52, val62", "val72", "val82"},
				"Header-Key":          []string{"val5, val6", "val7", "val8"},
				"Header-Key-Copy":     []string{"val5, val6", "val7", "val8"},
			},
			prefix: "Prefix",
			key:    "Header-Key",
			values: []string{"val1, val2", "val3", "val4"},
			expected: http.Header{
				"Header-Key":          []string{"val1, val2", "val3", "val4"},
				"Prefix-Header-Key":   []string{"val5, val6", "val7", "val8"},
				"Prefix-1-Header-Key": []string{"val50, val60", "val70", "val80"},
				"Prefix-2-Header-Key": []string{"val51, val61", "val71", "val81"},
				"Prefix-3-Header-Key": []string{"val52, val62", "val72", "val82"},
				"Header-Key-Copy":     []string{"val5, val6", "val7", "val8"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ReplaceHeader(tc.headers, tc.prefix, tc.key, tc.values...)

			if len(tc.expected) != len(tc.headers) {
				t.Fatalf("expected %v; got %v", tc.expected, tc.headers)
			}

			for key, values := range tc.expected {
				if len(values) != len(tc.headers[key]) {
					t.Fatalf("expected %v; got %v", tc.expected, tc.headers)
				}

				for i, value := range values {
					if value != tc.headers[key][i] {
						t.Fatalf("expected %v; got %v", tc.expected, tc.headers)
					}
				}
			}
		})
	}
}
