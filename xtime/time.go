// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xtime provides additional primitives and structures for measuring and displaying time.
package xtime

import (
	"errors"
	"strconv"
	"time"
)

const (
	hourInDay  = 24
	minInHour  = 60
	msecInSec  = 1e3
	nsecInMsec = 1e6
	secInMin   = 60
)

// A TimeMilli represents an instant in time with nanosecond precision,
// except for JSON/Text encoding/decoding which is of millisecond precision:
// - encoding uses xtime.RFC3339Milli layout,
// - decoding handles both Unix timestamps in milliseconds and xtime.RFC3339Milli layout.
//
// See time.Time for more information.
type TimeMilli struct {
	time.Time
}

// DateMilli returns the Time corresponding to
//	yyyy-mm-dd hh:mm:ss + msec milliseconds
// in the appropriate zone for that time in the given location.
//
// See time.Date for more information.
func DateMilli(year int, month time.Month, day, hour, min, sec, msec int, loc *time.Location) TimeMilli {
	// Normalize msec, sec, min, hour, overflowing into day.
	sec, msec = norm(sec, msec, msecInSec)
	min, sec = norm(min, sec, secInMin)
	hour, min = norm(hour, min, minInHour)
	day, hour = norm(day, hour, hourInDay)

	return TimeMilli{time.Date(year, month, day, hour, min, sec, msec*nsecInMsec, loc)}
}

// NowMilli returns the current local time.
//
// See time.Now for more information.
func NowMilli() TimeMilli {
	return TimeMilli{time.Now()}
}

// ToMilli is a convenience function to convert a time.Time into TimeMilli.
func ToMilli(t time.Time) TimeMilli {
	return TimeMilli{t}
}

// UnixMilli returns the local TimeMilli corresponding to the given Unix time,
// sec seconds and msec milliseconds since January 1, 1970 UTC.
// It is valid to pass msec outside the range [0, 999].
//
// See time.Unix for more information.
func UnixMilli(sec, msec int64) TimeMilli {
	if msec < 0 || msec >= msecInSec {
		n := msec / msecInSec
		sec += n
		msec -= n * msecInSec
		if msec < 0 {
			msec += msecInSec
			sec--
		}
	}
	return TimeMilli{time.Unix(sec, msec*nsecInMsec)}
}

// Add returns the time t+d.
//
// See time.Time.Add for more information.
func (t TimeMilli) Add(d time.Duration) TimeMilli {
	return TimeMilli{t.Time.Add(d)}
}

// AddDate returns the time corresponding to adding the
// given number of years, months, and days to t.
//
// See time.Time.AddDate for more information.
func (t TimeMilli) AddDate(years, months, days int) TimeMilli {
	return TimeMilli{t.Time.AddDate(years, months, days)}
}

// In returns a copy of t representing the same timestamp instant, but
// with the copy's location information set to loc for display
// purposes.
//
// See time.Time.In for more information.
func (t TimeMilli) In(loc *time.Location) TimeMilli {
	return TimeMilli{t.Time.In(loc)}
}

// Local returns t with the location set to local time.
//
// See time.Time.Local for more information.
func (t TimeMilli) Local() TimeMilli {
	return TimeMilli{t.Time.Local()}
}

// MarshalJSON implements the json.Marshaler interface.
// The time is a quoted string in RFC 3339 format, with sub-second precision added if present.
func (t TimeMilli) MarshalJSON() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		// RFC 3339 is clear that years are 4 digits exactly.
		// See golang.org/issue/4556#c15 for more discussion.
		return nil, errors.New("TimeMilli.MarshalJSON: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(RFC3339Milli)+2) //nolint:gomnd // Extra 2 double quotes.
	b = append(b, '"')
	b = t.AppendFormat(b, RFC3339Milli)
	b = append(b, '"')
	return b, nil
}

// MarshalText implements the encoding.TextMarshaler interface.
// The time is formatted in RFC 3339 format, with sub-second precision added if present.
func (t TimeMilli) MarshalText() ([]byte, error) {
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("TimeMilli.MarshalText: year outside of range [0,9999]")
	}

	b := make([]byte, 0, len(RFC3339Milli))
	return t.AppendFormat(b, RFC3339Milli), nil
}

// Millisecond returns the millisecond offset within the second specified by t,
// in the range [0, 999].
func (t TimeMilli) Millisecond() int {
	return t.Nanosecond() / nsecInMsec
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
//
// See time.Time.Round for more information.
func (t TimeMilli) Round(d time.Duration) TimeMilli {
	return TimeMilli{t.Time.Round(d)}
}

// T is a convenience method to access the underlying time.Time structure
// for compatibility with the Go standard time package.
func (t TimeMilli) T() time.Time {
	return t.Time
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
//
// See time.Time.Truncate for more information.
func (t TimeMilli) Truncate(d time.Duration) TimeMilli {
	return TimeMilli{t.Time.Truncate(d)}
}

// UTC returns t with the location set to UTC.
//
// See time.Time.UTC for more information.
func (t TimeMilli) UTC() TimeMilli {
	return TimeMilli{t.Time.UTC()}
}

// UnixMilli returns t as a Unix timestamp, the number of milliseconds elapsed
// since Time 1, 1970 UTC. The result does not depend on the location associated with it.
func (t TimeMilli) UnixMilli() int64 {
	return t.UnixNano() / nsecInMsec
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be:
// - either a quoted string in RFC 3339 format,
// - or a timestamp with millisecond precision expressed either as a number or a quoted string.
//
// See time.Time.UnmarshalJSON for more information.
func (t *TimeMilli) UnmarshalJSON(data []byte) error {
	b, e := 0, len(data)-1
	if len(data) > 1 && data[b] == '"' && data[e] == '"' {
		b++
		e--
	}

	if i, err := strconv.ParseInt(string(data[b:e+1]), 10, 64); err == nil { //nolint:gomnd // Decimal (base-10) integer.
		*t = UnixMilli(0, i)
		return nil
	}

	return t.Time.UnmarshalJSON(data)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The time is expected to be:
// - either a timestamp with millisecond precision,
// - or in RFC 3339 format.
//
// See time.Time.UnmarshalText for more information.
func (t *TimeMilli) UnmarshalText(data []byte) error {
	if i, err := strconv.ParseInt(string(data), 10, 64); err == nil { //nolint:gomnd // Decimal (base-10) integer.
		*t = UnixMilli(0, i)
		return nil
	}
	return t.Time.UnmarshalText(data)
}

// norm returns nhi, nlo such that
//	hi * base + lo == nhi * base + nlo
//	0 <= nlo < base
func norm(hi, lo, base int) (nhi, nlo int) {
	if lo < 0 {
		n := (-lo-1)/base + 1
		hi -= n
		lo += n * base
	}
	if lo >= base {
		n := lo / base
		hi += n
		lo -= n * base
	}
	return hi, lo
}
