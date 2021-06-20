package xtime_test

import (
	"testing"
	"time"

	. "github.com/jlourenc/xgo/xtime"
)

func TestDateStampMilli(t *testing.T) {
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
			got := DateStampMilli(tc.year, tc.month, tc.day, tc.hour, tc.min, tc.sec, tc.msec, tc.loc)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestNowStampMilli(t *testing.T) {
	before := time.Now()
	time.Sleep(time.Millisecond)
	got := NowStampMilli()
	time.Sleep(time.Millisecond)
	after := time.Now()

	if !got.After(before) || !after.After(got.T()) {
		t.Errorf("%s expected to be in range [%s, %s]", got, before, after)
	}
}

func TestToStampMilli(t *testing.T) {
	expected := time.Now()
	got := ToStampMilli(expected)

	if !got.Equal(expected) {
		t.Errorf("expected %s; got %s", expected, got)
	}
}

func TestUnixStampMilli(t *testing.T) {
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
			got := UnixStampMilli(tc.sec, tc.msec)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_Add(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp TimestampMilli
		duration  time.Duration
		expected  TimestampMilli
	}{
		{
			name:      "zero duration",
			timestamp: UnixStampMilli(1468181520, 499),
			duration:  0,
			expected:  UnixStampMilli(1468181520, 499),
		},
		{
			name:      "positive duration",
			timestamp: UnixStampMilli(1468181520, 499),
			duration:  20 * time.Second,
			expected:  UnixStampMilli(1468181540, 499),
		},
		{
			name:      "negative duration",
			timestamp: UnixStampMilli(1468181520, 499),
			duration:  -20 * time.Second,
			expected:  UnixStampMilli(1468181500, 499),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timestamp.Add(tc.duration)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_AddDate(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp TimestampMilli
		years     int
		months    int
		days      int
		expected  TimestampMilli
	}{
		{
			name:      "zero values",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			years:     0,
			months:    0,
			days:      0,
			expected:  DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
		{
			name:      "no overflow values",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			years:     3,
			months:    2,
			days:      1,
			expected:  DateStampMilli(2019, time.September, 11, 21, 12, 0, 499, time.UTC),
		},
		{
			name:      "with overflow values",
			timestamp: DateStampMilli(2016, time.December, 31, 21, 12, 0, 499, time.UTC),
			years:     1,
			months:    1,
			days:      1,
			expected:  DateStampMilli(2018, time.February, 1, 21, 12, 0, 499, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timestamp.AddDate(tc.years, tc.months, tc.days)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_In(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp TimestampMilli
		location  *time.Location
		expected  TimestampMilli
	}{
		{
			name:      "UTC to UTC",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			location:  time.UTC,
			expected:  DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
		{
			name:      "UTC to CET",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			location:  time.FixedZone("CET", 2*60*60),
			expected:  DateStampMilli(2016, time.July, 10, 23, 12, 0, 499, time.FixedZone("CET", 2*60*60)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timestamp.In(tc.location)
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_Local(t *testing.T) {
	_, localOffset := time.Now().Local().Zone()
	testCases := []struct {
		name      string
		timestamp TimestampMilli
		expected  TimestampMilli
	}{
		{
			name:      "from UTC",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expected:  DateStampMilli(2016, time.July, 10, 21, 12, localOffset, 499, time.Local),
		},
		{
			name:      "from Local",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.Local),
			expected:  DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.Local),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timestamp.Local()
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name          string
		timestamp     TimestampMilli
		expectedBytes []byte
		expectedErr   error
	}{
		{
			name:          "UTC - with msec",
			timestamp:     DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedBytes: []byte(`1468185120499`),
			expectedErr:   nil,
		},
		{
			name:          "zone info - no msec",
			timestamp:     DateStampMilli(2016, time.July, 10, 21, 12, 0, 0, time.FixedZone("CET", 2*60*60)),
			expectedBytes: []byte(`1468177920000`),
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBytes, gotErr := tc.timestamp.MarshalJSON()

			if string(tc.expectedBytes) != string(gotBytes) {
				t.Errorf("expected bytes %s; got %s", tc.expectedBytes, gotBytes)
			}

			if (tc.expectedErr == nil && gotErr != nil) ||
				(tc.expectedErr != nil && tc.expectedErr.Error() != gotErr.Error()) {
				t.Errorf("expected error %s; got %s", tc.expectedErr, gotErr)
			}
		})
	}
}

func TestTimestampMilli_MarshalText(t *testing.T) {
	testCases := []struct {
		name          string
		timestamp     TimestampMilli
		expectedBytes []byte
		expectedErr   error
	}{
		{
			name:          "UTC - with msec",
			timestamp:     DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedBytes: []byte(`1468185120499`),
			expectedErr:   nil,
		},
		{
			name:          "zone info - no msec",
			timestamp:     DateStampMilli(2016, time.July, 10, 21, 12, 0, 0, time.FixedZone("CET", 2*60*60)),
			expectedBytes: []byte(`1468177920000`),
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotBytes, gotErr := tc.timestamp.MarshalText()

			if string(tc.expectedBytes) != string(gotBytes) {
				t.Errorf("expected bytes %s; got %s", tc.expectedBytes, gotBytes)
			}

			if (tc.expectedErr == nil && gotErr != nil) ||
				(tc.expectedErr != nil && tc.expectedErr.Error() != gotErr.Error()) {
				t.Errorf("expected error %s; got %s", tc.expectedErr, gotErr)
			}
		})
	}
}

func TestTimestampMilli_Millisecond(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp TimestampMilli
		expected  int
	}{
		{
			name:      "no msec",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 0, time.UTC),
			expected:  0,
		},
		{
			name:      "with msec",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expected:  499,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timestamp.Millisecond()
			if tc.expected != got {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_Round(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp TimestampMilli
		duration  time.Duration
		expected  TimestampMilli
	}{
		{
			name:      "nearest second",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC),
			duration:  2 * time.Second,
			expected:  DateStampMilli(2016, time.July, 10, 21, 12, 2, 0, time.UTC),
		},
		{
			name:      "nearest hour",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			duration:  time.Hour,
			expected:  DateStampMilli(2016, time.July, 10, 21, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timestamp.Round(tc.duration)
			if tc.expected != got {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_T(t *testing.T) {
	x := DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	expected := time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.UTC)
	got := x.T()

	if !expected.Equal(got) {
		t.Errorf("expected %s; got %s", expected, got)
	}
}

func TestTimestampMilli_Truncate(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp TimestampMilli
		duration  time.Duration
		expected  TimestampMilli
	}{
		{
			name:      "rounding down to a multiple of 2 seconds",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC),
			duration:  2 * time.Second,
			expected:  DateStampMilli(2016, time.July, 10, 21, 12, 0, 0, time.UTC),
		},
		{
			name:      "rounding down to nearest hour",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			duration:  time.Hour,
			expected:  DateStampMilli(2016, time.July, 10, 21, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timestamp.Truncate(tc.duration)
			if tc.expected != got {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_UTC(t *testing.T) {
	_, localOffset := time.Now().Local().Zone()
	testCases := []struct {
		name      string
		timestamp TimestampMilli
		expected  TimestampMilli
	}{
		{
			name:      "from UTC",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expected:  DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
		},
		{
			name:      "from Local",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.Local),
			expected:  DateStampMilli(2016, time.July, 10, 21, 12, -localOffset, 499, time.UTC),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timestamp.UTC()
			if !tc.expected.Equal(got.T()) {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_UnixMilli(t *testing.T) {
	testCases := []struct {
		name      string
		timestamp TimestampMilli
		expected  int64
	}{
		{
			name:      "no msec",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 0, time.UTC),
			expected:  1468185120000,
		},
		{
			name:      "with msec",
			timestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expected:  1468185120499,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.timestamp.UnixMilli()
			if tc.expected != got {
				t.Errorf("expected %d; got %d", tc.expected, got)
			}
		})
	}
}

func TestTimestampMilli_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name              string
		data              []byte
		expectedTimestamp TimestampMilli
		expectedErr       error
	}{
		{
			name:              "double-quoted string timestamp",
			data:              []byte(`"1468185120499"`),
			expectedTimestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedErr:       nil,
		},
		{
			name:              "number timestamp",
			data:              []byte(`1468185120499`),
			expectedTimestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedErr:       nil,
		},
		{
			name:              "RFC3339 string",
			data:              []byte(`"2016-07-10T21:12:00+02:00"`),
			expectedTimestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 0, time.FixedZone("CET", 2*60*60)),
		},
		{
			name:              "RFC3339Milli string",
			data:              []byte(`"2016-07-10T21:12:00.499+02:00"`),
			expectedTimestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.FixedZone("CET", 2*60*60)),
		},
		{
			name:              "RFC3339Nano string",
			data:              []byte(`"2016-07-10T21:12:00.499+02:00"`),
			expectedTimestamp: ToStampMilli(time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.FixedZone("CET", 2*60*60))),
		},
	}

	for _, tc := range testCases {
		var gotTime TimestampMilli
		gotErr := gotTime.UnmarshalJSON(tc.data)

		if !tc.expectedTimestamp.Equal(gotTime.T()) {
			t.Errorf("expected time %s; got %s", tc.expectedTimestamp, gotTime)
		}

		if (tc.expectedErr == nil && gotErr != nil) ||
			(tc.expectedErr != nil && tc.expectedErr.Error() != gotErr.Error()) {
			t.Errorf("expected error %s; got %s", tc.expectedErr, gotErr)
		}
	}
}

func TestTimestampMilli_UnmarshalText(t *testing.T) {
	testCases := []struct {
		name              string
		data              []byte
		expectedTimestamp TimestampMilli
		expectedErr       error
	}{
		{
			name:              "timestamp",
			data:              []byte(`1468185120499`),
			expectedTimestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC),
			expectedErr:       nil,
		},
		{
			name:              "RFC3339",
			data:              []byte(`2016-07-10T21:12:00+02:00`),
			expectedTimestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 0, time.FixedZone("CET", 2*60*60)),
		},
		{
			name:              "RFC3339Milli",
			data:              []byte(`2016-07-10T21:12:00.499+02:00`),
			expectedTimestamp: DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.FixedZone("CET", 2*60*60)),
		},
		{
			name:              "RFC3339Nano",
			data:              []byte(`2016-07-10T21:12:00.499+02:00`),
			expectedTimestamp: ToStampMilli(time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.FixedZone("CET", 2*60*60))),
		},
	}

	for _, tc := range testCases {
		var gotTime TimestampMilli
		gotErr := gotTime.UnmarshalText(tc.data)

		if !tc.expectedTimestamp.Equal(gotTime.T()) {
			t.Errorf("expected time %s; got %s", tc.expectedTimestamp, gotTime)
		}

		if (tc.expectedErr == nil && gotErr != nil) ||
			(tc.expectedErr != nil && tc.expectedErr.Error() != gotErr.Error()) {
			t.Errorf("expected error %s; got %s", tc.expectedErr, gotErr)
		}
	}
}
