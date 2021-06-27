// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xtime

import (
	"strconv"
	"time"
)

// A TimestampMilli represents an instant in time with nanosecond precision,
// except for JSON/Text encoding/decoding which is of millisecond precision:
// 1) encoding uses Unix timestamps in milliseconds,
// 2) decoding handles both Unix timestamps in milliseconds and xtime.RFC3339Milli layout.
//
// See time.Time for more information.
type TimestampMilli struct {
	time.Time
}

// DateStampMilli returns the TimestampMilli corresponding to
//	yyyy-mm-dd hh:mm:ss + msec milliseconds
// in the appropriate zone for that timestamp in the given location.
//
// See time.Date for more information.
func DateStampMilli(year int, month time.Month, day, hour, min, sec, msec int, loc *time.Location) TimestampMilli {
	// Normalize msec, sec, min, hour, overflowing into day.
	sec, msec = norm(sec, msec, msecsInSec)
	min, sec = norm(min, sec, secsInMin)
	hour, min = norm(hour, min, minsInHour)
	day, hour = norm(day, hour, hoursInDay)

	return TimestampMilli{time.Date(year, month, day, hour, min, sec, msec*nsecsInMsec, loc)}
}

// NowStampMilli returns the current local time as TimestampMilli.
//
// See time.Now for more information.
func NowStampMilli() TimestampMilli {
	return TimestampMilli{time.Now()}
}

// ToStampMilli is a convenience function to convert a time.Time into TimestampMilli.
func ToStampMilli(t time.Time) TimestampMilli {
	return TimestampMilli{Time: t}
}

// UnixStampMilli returns the local TimestampMilli corresponding to the given Unix time,
// sec seconds and msec milliseconds since January 1, 1970 UTC.
// It is valid to pass msec outside the range [0, 999].
//
// See time.Unix for more information.
func UnixStampMilli(sec, msec int64) TimestampMilli {
	if msec < 0 || msec >= msecsInSec {
		n := msec / msecsInSec
		sec += n
		msec -= n * msecsInSec
		if msec < 0 {
			msec += msecsInSec
			sec--
		}
	}
	return TimestampMilli{time.Unix(sec, msec*nsecsInMsec)}
}

// Add returns the time t+d.
//
// See time.Time.Add for more information.
func (t TimestampMilli) Add(d time.Duration) TimestampMilli {
	return TimestampMilli{t.Time.Add(d)}
}

// AddDate returns the time corresponding to adding the
// given number of years, months, and days to t.
//
// See time.Time.AddDate for more information.
func (t TimestampMilli) AddDate(years, months, days int) TimestampMilli {
	return TimestampMilli{t.Time.AddDate(years, months, days)}
}

// In returns a copy of t representing the same timestamp instant, but
// with the copy's location information set to loc for display
// purposes.
//
// See time.Time.In for more information.
func (t TimestampMilli) In(loc *time.Location) TimestampMilli {
	return TimestampMilli{t.Time.In(loc)}
}

// Local returns t with the location set to local time.
//
// See time.Time.Local for more information.
func (t TimestampMilli) Local() TimestampMilli {
	return TimestampMilli{t.Time.Local()}
}

// MarshalJSON implements the json.Marshaler interface.
// The timestamp is a Unix timestamp with millisecond precision.
func (t TimestampMilli) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.UnixMilli(), 10)), nil //nolint:gomnd // Decimal (base-10) integer.
}

// MarshalText implements the encoding.TextMarshaler interface.
// The timestamp is a Unix timestamp with millisecond precision.
func (t TimestampMilli) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatInt(t.UnixMilli(), 10)), nil //nolint:gomnd // Decimal (base-10) integer.
}

// Millisecond returns the millisecond offset within the second specified by t,
// in the range [0, 999].
func (t TimestampMilli) Millisecond() int {
	return t.Nanosecond() / nsecsInMsec
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
//
// See time.Time.Round for more information.
func (t TimestampMilli) Round(d time.Duration) TimestampMilli {
	return TimestampMilli{t.Time.Round(d)}
}

// T is a convenience method to access the underlying time.Time structure
// for compatibility with the Go standard time package.
func (t TimestampMilli) T() time.Time {
	return t.Time
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
//
// See time.Time.Truncate for more information.
func (t TimestampMilli) Truncate(d time.Duration) TimestampMilli {
	return TimestampMilli{t.Time.Truncate(d)}
}

// UTC returns t with the location set to UTC.
//
// See time.Time.UTC for more information.
func (t TimestampMilli) UTC() TimestampMilli {
	return TimestampMilli{t.Time.UTC()}
}

// UnixMilli returns t as a Unix timestamp, the number of milliseconds elapsed
// since Time 1, 1970 UTC. The result does not depend on the location associated with it.
func (t TimestampMilli) UnixMilli() int64 {
	return t.UnixNano() / nsecsInMsec
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be either
// 1) a quoted string in RFC 3339 format, or
// 2) a timestamp with millisecond precision expressed either as a number or a quoted string.
//
// See time.Time.UnmarshalJSON for more information.
func (t *TimestampMilli) UnmarshalJSON(data []byte) error {
	b, e := 0, len(data)-1
	if len(data) > 1 && data[b] == '"' && data[e] == '"' {
		b++
		e--
	}

	if i, err := strconv.ParseInt(string(data[b:e+1]), 10, 64); err == nil { //nolint:gomnd // Decimal (base-10) integer.
		*t = UnixStampMilli(0, i)
		return nil
	}

	return t.Time.UnmarshalJSON(data)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be either
// 1) a timestamp with millisecond precision, or
// 2) in RFC 3339 format.
//
// See time.Time.UnmarshalText for more information.
func (t *TimestampMilli) UnmarshalText(data []byte) error {
	if i, err := strconv.ParseInt(string(data), 10, 64); err == nil { //nolint:gomnd // Decimal (base-10) integer.
		*t = UnixStampMilli(0, i)
		return nil
	}
	return t.Time.UnmarshalText(data)
}
