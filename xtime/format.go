// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xtime

import (
	"time"
)

// These are predefined extra layouts to use in time.Format and time.Parse.
const (
	RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"
)

// ParseMilli parses a formatted string and returns the time value it represents as TimeMilli.
//
// See time.Parse for more information.
func ParseMilli(layout, value string) (TimeMilli, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return TimeMilli{}, err
	}
	return TimeMilli{t}, nil
}

// ParseMilliInLocation is like ParseMilli but with an extra time.Location argument.
//
// See time.ParseInLocation for more information.
func ParseMilliInLocation(layout, value string, loc *time.Location) (TimeMilli, error) {
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return TimeMilli{}, err
	}
	return TimeMilli{t}, nil
}

// ParseStampMilli parses a formatted string and returns the time value it represents as TimestampMilli.
//
// See time.Parse for more information.
func ParseStampMilli(layout, value string) (TimestampMilli, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return TimestampMilli{}, err
	}
	return TimestampMilli{t}, nil
}

// ParseStampMilliInLocation is like ParseStampMilli but with an extra time.Location argument.
//
// See time.ParseInLocation for more information.
func ParseStampMilliInLocation(layout, value string, loc *time.Location) (TimestampMilli, error) {
	t, err := time.ParseInLocation(layout, value, loc)
	if err != nil {
		return TimestampMilli{}, err
	}
	return TimestampMilli{t}, nil
}
