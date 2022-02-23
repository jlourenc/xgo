// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xtime_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	. "github.com/jlourenc/xgo/xtime"
)

func TestDateMilli(t *testing.T) {
	testCases := []struct {
		name     string
		year     int
		month    time.Month
		day      int
		hour     int
		min      int
		sec      int
		msec     int
		loc      *time.Location
		expected time.Time
	}{
		{
			name:     "UTC - no overflow",
			year:     2016,
			month:    time.July,
			day:      10,
			hour:     21,
			min:      12,
			sec:      0,
			msec:     499,
			loc:      time.UTC,
			expected: time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.UTC),
		},
		{
			name:     "Local - no overflow",
			year:     2016,
			month:    time.July,
			day:      10,
			hour:     21,
			min:      12,
			sec:      0,
			msec:     499,
			loc:      time.Local,
			expected: time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.Local),
		},
		{
			name:     "UTC - positive overflow",
			year:     2016,
			month:    time.December,
			day:      31,
			hour:     23,
			min:      59,
			sec:      59,
			msec:     1499,
			loc:      time.UTC,
			expected: time.Date(2017, time.January, 1, 0, 0, 0, 499000000, time.UTC),
		},
		{
			name:     "UTC - negative overflow",
			year:     2016,
			month:    time.July,
			day:      10,
			hour:     21,
			min:      12,
			sec:      0,
			msec:     -1,
			loc:      time.UTC,
			expected: time.Date(2016, time.July, 10, 21, 11, 59, 999000000, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := DateMilli(tc.year, tc.month, tc.day, tc.hour, tc.min, tc.sec, tc.msec, tc.loc)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestNowMilli(t *testing.T) {
	before := time.Now()
	time.Sleep(time.Millisecond)
	got := NowMilli()
	time.Sleep(time.Millisecond)
	after := time.Now()

	if !got.After(before) || !after.After(got.T()) {
		t.Errorf("%s expected to be in range [%s, %s]", got, before, after)
	}
}

func TestToMilli(t *testing.T) {
	expected := time.Now()
	got := ToMilli(expected)

	if !got.Equal(expected) {
		t.Errorf("expected %s; got %s", expected, got)
	}
}

func TestUnixMilli(t *testing.T) {
	testCases := []struct {
		name     string
		sec      int64
		msec     int64
		expected time.Time
	}{
		{
			name:     "within msec range",
			sec:      1468181520,
			msec:     499,
			expected: time.Unix(1468181520, 499000000),
		},
		{
			name:     "outside msec range - positive",
			sec:      1468181520,
			msec:     61499,
			expected: time.Unix(1468181581, 499000000),
		},
		{
			name:     "outside msec range - negative",
			sec:      1468181520,
			msec:     -1,
			expected: time.Unix(1468181519, 999000000),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := UnixMilli(tc.sec, tc.msec)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_Add(t *testing.T) {
	testCases := []struct {
		name     string
		time     TimeMilli
		duration time.Duration
		expected TimeMilli
	}{
		{
			name:     "zero duration",
			time:     UnixMilli(1468181520, 499),
			duration: 0,
			expected: UnixMilli(1468181520, 499),
		},
		{
			name:     "positive duration",
			time:     UnixMilli(1468181520, 499),
			duration: 20 * time.Second,
			expected: UnixMilli(1468181540, 499),
		},
		{
			name:     "negative duration",
			time:     UnixMilli(1468181520, 499),
			duration: -20 * time.Second,
			expected: UnixMilli(1468181500, 499),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.time.Add(tc.duration)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_AddDate(t *testing.T) {
	testCases := []struct {
		name     string
		time     TimeMilli
		years    int
		months   int
		days     int
		expected TimeMilli
	}{
		{
			name:     "zero values",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			years:    0,
			months:   0,
			days:     0,
			expected: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
		{
			name:     "no overflow values",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			years:    3,
			months:   2,
			days:     1,
			expected: DateMilli(2019, time.September, 11, 21, 12, 0, 499, time.UTC),
		},
		{
			name:     "with overflow values",
			time:     DateMilli(2016, time.December, 31, 21, 12, 0, 499, time.UTC),
			years:    1,
			months:   1,
			days:     1,
			expected: DateMilli(2018, time.February, 1, 21, 12, 0, 499, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.time.AddDate(tc.years, tc.months, tc.days)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_In(t *testing.T) {
	testCases := []struct {
		name     string
		time     TimeMilli
		location *time.Location
		expected TimeMilli
	}{
		{
			name:     "UTC to UTC",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			location: time.UTC,
			expected: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
		{
			name:     "UTC to CET",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			location: time.FixedZone("CET", 2*60*60),
			expected: DateMilli(2016, time.July, 10, 23, 12, 0, 499, time.FixedZone("CET", 2*60*60)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.time.In(tc.location)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_Local(t *testing.T) {
	_, localOffset := time.Now().Local().Zone()
	testCases := []struct {
		name     string
		time     TimeMilli
		expected TimeMilli
	}{
		{
			name:     "from UTC",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expected: DateMilli(2016, time.July, 10, 21, 12, localOffset, 499, time.FixedZone("local", localOffset)),
		},
		{
			name:     "from Local",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.Local),
			expected: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.Local),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.time.Local()
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name          string
		time          TimeMilli
		expectedBytes []byte
		expectedErr   error
	}{
		{
			name:          "invalid year",
			time:          DateMilli(10001, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedBytes: nil,
			expectedErr:   errors.New("TimeMilli.MarshalJSON: year outside of range [0,9999]"),
		},
		{
			name:          "UTC - with msec",
			time:          DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedBytes: []byte(`"2016-07-10T21:12:00.499Z"`),
			expectedErr:   nil,
		},
		{
			name:          "zone info - no msec",
			time:          DateMilli(2016, time.July, 10, 21, 12, 0, 0, time.FixedZone("CET", 2*60*60)),
			expectedBytes: []byte(`"2016-07-10T21:12:00+02:00"`),
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBytes, gotErr := tc.time.MarshalJSON()

			if !bytes.Equal(tc.expectedBytes, gotBytes) {
				t.Errorf("expected bytes %s; got %s", tc.expectedBytes, gotBytes)
			}

			if (tc.expectedErr == nil && gotErr != nil) ||
				(tc.expectedErr != nil && tc.expectedErr.Error() != gotErr.Error()) {
				t.Errorf("expected error %s; got %s", tc.expectedErr, gotErr)
			}
		})
	}
}

func TestTimeMilli_MarshalText(t *testing.T) {
	testCases := []struct {
		name          string
		time          TimeMilli
		expectedBytes []byte
		expectedErr   error
	}{
		{
			name:          "invalid year",
			time:          DateMilli(10001, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedBytes: nil,
			expectedErr:   errors.New("TimeMilli.MarshalText: year outside of range [0,9999]"),
		},
		{
			name:          "UTC - with msec",
			time:          DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedBytes: []byte(`2016-07-10T21:12:00.499Z`),
			expectedErr:   nil,
		},
		{
			name:          "zone info - no msec",
			time:          DateMilli(2016, time.July, 10, 21, 12, 0, 0, time.FixedZone("CET", 2*60*60)),
			expectedBytes: []byte(`2016-07-10T21:12:00+02:00`),
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBytes, gotErr := tc.time.MarshalText()

			if !bytes.Equal(tc.expectedBytes, gotBytes) {
				t.Errorf("expected bytes %s; got %s", tc.expectedBytes, gotBytes)
			}

			if (tc.expectedErr == nil && gotErr != nil) ||
				(tc.expectedErr != nil && tc.expectedErr.Error() != gotErr.Error()) {
				t.Errorf("expected error %s; got %s", tc.expectedErr, gotErr)
			}
		})
	}
}

func TestTimeMilli_Millisecond(t *testing.T) {
	testCases := []struct {
		name     string
		time     TimeMilli
		expected int
	}{
		{
			name:     "no msec",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "with msec",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expected: 499,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.time.Millisecond()
			if tc.expected != got {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_Round(t *testing.T) {
	testCases := []struct {
		name     string
		time     TimeMilli
		duration time.Duration
		expected TimeMilli
	}{
		{
			name:     "nearest second",
			time:     DateMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC),
			duration: 2 * time.Second,
			expected: DateMilli(2016, time.July, 10, 21, 12, 2, 0, time.UTC),
		},
		{
			name:     "nearest hour",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			duration: time.Hour,
			expected: DateMilli(2016, time.July, 10, 21, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.time.Round(tc.duration)
			if tc.expected != got {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_T(t *testing.T) {
	x := DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	expected := time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.UTC)
	got := x.T()

	if !expected.Equal(got) {
		t.Errorf("expected %s; got %s", expected, got)
	}
}

func TestTimeMilli_Truncate(t *testing.T) {
	testCases := []struct {
		name     string
		time     TimeMilli
		duration time.Duration
		expected TimeMilli
	}{
		{
			name:     "rounding down to a multiple of 2 seconds",
			time:     DateMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC),
			duration: 2 * time.Second,
			expected: DateMilli(2016, time.July, 10, 21, 12, 0, 0, time.UTC),
		},
		{
			name:     "rounding down to nearest hour",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			duration: time.Hour,
			expected: DateMilli(2016, time.July, 10, 21, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.time.Truncate(tc.duration)
			if tc.expected != got {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_UTC(t *testing.T) {
	_, localOffset := time.Now().Local().Zone()
	testCases := []struct {
		name     string
		time     TimeMilli
		expected TimeMilli
	}{
		{
			name:     "from UTC",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expected: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
		{
			name:     "from Local",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.FixedZone("local", localOffset)),
			expected: DateMilli(2016, time.July, 10, 21, 12, -localOffset, 499, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.time.UTC()
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_UnixMilli(t *testing.T) {
	testCases := []struct {
		name     string
		time     TimeMilli
		expected int64
	}{
		{
			name:     "no msec",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 0, time.UTC),
			expected: 1468185120000,
		},
		{
			name:     "with msec",
			time:     DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expected: 1468185120499,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.time.UnixMilli()
			if tc.expected != got {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestTimeMilli_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		data         []byte
		expectedTime TimeMilli
		expectedErr  error
	}{
		{
			name:         "double-quoted string timestamp",
			data:         []byte(`"1468185120499"`),
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedErr:  nil,
		},
		{
			name:         "number timestamp",
			data:         []byte(`1468185120499`),
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedErr:  nil,
		},
		{
			name:         "RFC3339 string",
			data:         []byte(`"2016-07-10T21:12:00+02:00"`),
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 0, time.FixedZone("CET", 2*60*60)),
		},
		{
			name:         "RFC3339Milli string",
			data:         []byte(`"2016-07-10T21:12:00.499+02:00"`),
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.FixedZone("CET", 2*60*60)),
		},
		{
			name:         "RFC3339Nano string",
			data:         []byte(`"2016-07-10T21:12:00.499+02:00"`),
			expectedTime: ToMilli(time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.FixedZone("CET", 2*60*60))),
		},
	}

	for _, tc := range testCases {
		var gotTime TimeMilli
		gotErr := gotTime.UnmarshalJSON(tc.data)

		if !tc.expectedTime.Equal(gotTime.T()) {
			t.Errorf("expected time %s; got %s", tc.expectedTime, gotTime)
		}

		if (tc.expectedErr == nil && gotErr != nil) ||
			(tc.expectedErr != nil && tc.expectedErr.Error() != gotErr.Error()) {
			t.Errorf("expected error %s; got %s", tc.expectedErr, gotErr)
		}
	}
}

func TestTimeMilli_UnmarshalText(t *testing.T) {
	testCases := []struct {
		name         string
		data         []byte
		expectedTime TimeMilli
		expectedErr  error
	}{
		{
			name:         "timestamp",
			data:         []byte(`1468185120499`),
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedErr:  nil,
		},
		{
			name:         "RFC3339",
			data:         []byte(`2016-07-10T21:12:00+02:00`),
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 0, time.FixedZone("CET", 2*60*60)),
		},
		{
			name:         "RFC3339Milli",
			data:         []byte(`2016-07-10T21:12:00.499+02:00`),
			expectedTime: DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.FixedZone("CET", 2*60*60)),
		},
		{
			name:         "RFC3339Nano",
			data:         []byte(`2016-07-10T21:12:00.499+02:00`),
			expectedTime: ToMilli(time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.FixedZone("CET", 2*60*60))),
		},
	}

	for _, tc := range testCases {
		var gotTime TimeMilli
		gotErr := gotTime.UnmarshalText(tc.data)

		if !tc.expectedTime.Equal(gotTime.T()) {
			t.Errorf("expected time %s; got %s", tc.expectedTime, gotTime)
		}

		if (tc.expectedErr == nil && gotErr != nil) ||
			(tc.expectedErr != nil && tc.expectedErr.Error() != gotErr.Error()) {
			t.Errorf("expected error %s; got %s", tc.expectedErr, gotErr)
		}
	}
}
