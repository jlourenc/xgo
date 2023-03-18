// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xtime_test

import (
	"testing"
	"time"

	. "github.com/jlourenc/xgo/xtime"
)

func TestLayouts(t *testing.T) {
	x := time.Date(2016, 7, 10, 21, 12, 0, 499999999, time.UTC)

	testCases := []struct {
		name   string
		layout string
		value  string
		offset time.Duration
	}{
		{
			name:   "RFC3339Milli",
			layout: RFC3339Milli,
			value:  "2016-07-10T21:12:00.499Z",
			offset: 999999 * time.Nanosecond,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt := x.Format(tc.layout)
			if fmt != tc.value {
				t.Errorf("expected %s; got %s", tc.value, fmt)
			}

			p, err := time.Parse(tc.layout, tc.value)
			if err != nil {
				t.Errorf("no error expected; got %s", err)
			} else if p.Add(tc.offset) != x {
				t.Errorf("expected %v; got %v", x, p.Add(tc.offset))
			}
		})
	}
}

func TestParseMilli(t *testing.T) {
	testCases := []struct {
		name         string
		layout       string
		value        string
		expectedTime TimeMilli
		expectedErr  bool
	}{
		{
			name:        "invalid value",
			layout:      RFC3339Milli,
			value:       "invalid",
			expectedErr: true,
		},
		{
			name:         "RFC3339Milli",
			layout:       RFC3339Milli,
			value:        "2016-07-10T21:12:00.499Z",
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x, err := ParseMilli(tc.layout, tc.value)

			if tc.expectedErr && err == nil {
				t.Error("error expected; got nil")
			}
			if !tc.expectedErr && err != nil {
				t.Errorf("no error expected; got %s", err)
			}

			if !tc.expectedTime.Equal(x.T()) {
				t.Errorf("expected %s; got %s", tc.expectedTime, x)
			}
		})
	}
}

func TestParseMilliInLocation(t *testing.T) {
	testCases := []struct {
		name         string
		layout       string
		value        string
		location     *time.Location
		expectedTime TimeMilli
		expectedErr  bool
	}{
		{
			name:        "invalid value",
			layout:      RFC3339Milli,
			value:       "invalid",
			expectedErr: true,
		},
		{
			name:         "RFC3339Milli - nil location",
			layout:       RFC3339Milli,
			value:        "2016-07-10T21:12:00.499Z",
			location:     nil,
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
		{
			name:         "RFC3339Milli - CET",
			layout:       RFC3339Milli,
			value:        "2016-07-10T21:12:00.499+02:00",
			location:     time.FixedZone("CET", 2*60*60),
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.FixedZone("CET", 2*60*60)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x, err := ParseMilliInLocation(tc.layout, tc.value, tc.location)

			if tc.expectedErr && err == nil {
				t.Error("error expected; got nil")
			}
			if !tc.expectedErr && err != nil {
				t.Errorf("no error expected; got %s", err)
			}

			if !tc.expectedTime.Equal(x.T()) {
				t.Errorf("expected %s; got %s", tc.expectedTime, x)
			}
		})
	}
}

func TestParseStampMilli(t *testing.T) {
	testCases := []struct {
		name         string
		layout       string
		value        string
		expectedTime TimestampMilli
		expectedErr  bool
	}{
		{
			name:        "invalid value",
			layout:      RFC3339Milli,
			value:       "invalid",
			expectedErr: true,
		},
		{
			name:         "RFC3339Milli",
			layout:       RFC3339Milli,
			value:        "2016-07-10T21:12:00.499Z",
			expectedTime: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x, err := ParseStampMilli(tc.layout, tc.value)

			if tc.expectedErr && err == nil {
				t.Error("error expected; got nil")
			}
			if !tc.expectedErr && err != nil {
				t.Errorf("no error expected; got %s", err)
			}

			if !tc.expectedTime.Equal(x.T()) {
				t.Errorf("expected %s; got %s", tc.expectedTime, x)
			}
		})
	}
}

func TestParseStampMilliInLocation(t *testing.T) {
	testCases := []struct {
		name         string
		layout       string
		value        string
		location     *time.Location
		expectedTime TimestampMilli
		expectedErr  bool
	}{
		{
			name:        "invalid value",
			layout:      RFC3339Milli,
			value:       "invalid",
			expectedErr: true,
		},
		{
			name:         "RFC3339Milli - nil location",
			layout:       RFC3339Milli,
			value:        "2016-07-10T21:12:00.499Z",
			location:     nil,
			expectedTime: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
		{
			name:         "RFC3339Milli - CET",
			layout:       RFC3339Milli,
			value:        "2016-07-10T21:12:00.499+02:00",
			location:     time.FixedZone("CET", 2*60*60),
			expectedTime: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.FixedZone("CET", 2*60*60)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			x, err := ParseStampMilliInLocation(tc.layout, tc.value, tc.location)

			if tc.expectedErr && err == nil {
				t.Error("error expected; got nil")
			}
			if !tc.expectedErr && err != nil {
				t.Errorf("no error expected; got %s", err)
			}

			if !tc.expectedTime.Equal(x.T()) {
				t.Errorf("expected %s; got %s", tc.expectedTime, x)
			}
		})
	}
}
