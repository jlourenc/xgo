// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xtime_test

import (
	"fmt"
	"time"

	"github.com/jlourenc/xgo/xtime"
)

func ExampleDateMilli() {
	t := xtime.DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	fmt.Printf("Go launched at %s\n", t)
	// Output: Go launched at 2016-07-10 21:12:00.499 +0000 UTC
}

func ExampleParseMilli() {
	t, err := xtime.ParseMilli(xtime.RFC3339Milli, "2016-07-10T21:12:00.499Z")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", t)
	// Output: 2016-07-10 21:12:00.499 +0000 UTC
}

func ExampleParseMilliInLocation() {
	t, err := xtime.ParseMilliInLocation(xtime.RFC3339Milli, "2016-07-10T21:12:00.499+02:00", time.FixedZone("CET", 2*60*60))
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", t)
	// Output: 2016-07-10 21:12:00.499 +0200 CET
}

func ExampleToMilli() {
	t := xtime.ToMilli(time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.UTC))
	fmt.Printf("Go launched at %s\n", t)
	// Output: Go launched at 2016-07-10 21:12:00.499 +0000 UTC
}

func ExampleTimeMilli_Add() {
	start := xtime.DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	afterTenSeconds := start.Add(time.Second * 10)
	afterTenMinutes := start.Add(time.Minute * 10)
	afterTenHours := start.Add(time.Hour * 10)
	afterTenDays := start.Add(time.Hour * 24 * 10)

	fmt.Printf("start = %v\n", start)
	fmt.Printf("start.Add(time.Second * 10) = %s\n", afterTenSeconds)
	fmt.Printf("start.Add(time.Minute * 10) = %s\n", afterTenMinutes)
	fmt.Printf("start.Add(time.Hour * 10) = %s\n", afterTenHours)
	fmt.Printf("start.Add(time.Hour * 24 * 10) = %s\n", afterTenDays)

	// Output:
	// start = 2016-07-10 21:12:00.499 +0000 UTC
	// start.Add(time.Second * 10) = 2016-07-10 21:12:10.499 +0000 UTC
	// start.Add(time.Minute * 10) = 2016-07-10 21:22:00.499 +0000 UTC
	// start.Add(time.Hour * 10) = 2016-07-11 07:12:00.499 +0000 UTC
	// start.Add(time.Hour * 24 * 10) = 2016-07-20 21:12:00.499 +0000 UTC
}

func ExampleTimeMilli_AddDate() {
	start := xtime.DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	oneDayLater := start.AddDate(0, 0, 1)
	oneMonthLater := start.AddDate(0, 1, 0)
	oneYearLater := start.AddDate(1, 0, 0)

	fmt.Printf("oneDayLater: start.AddDate(0, 0, 1) = %s\n", oneDayLater)
	fmt.Printf("oneMonthLater: start.AddDate(0, 1, 0) = %s\n", oneMonthLater)
	fmt.Printf("oneYearLater: start.AddDate(1, 0, 0) = %s\n", oneYearLater)

	// Output:
	// oneDayLater: start.AddDate(0, 0, 1) = 2016-07-11 21:12:00.499 +0000 UTC
	// oneMonthLater: start.AddDate(0, 1, 0) = 2016-08-10 21:12:00.499 +0000 UTC
	// oneYearLater: start.AddDate(1, 0, 0) = 2017-07-10 21:12:00.499 +0000 UTC
}

func ExampleTimeMilli_In() {
	t := xtime.DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	t = t.In(time.FixedZone("CET", 2*60*60))
	fmt.Printf("%s\n", t)
	// Output: 2016-07-10 23:12:00.499 +0200 CET
}

func ExampleTimeMilli_Millisecond() {
	t := xtime.DateMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	msec := t.Millisecond()
	fmt.Printf("msec: %d\n", msec)
	// Output: msec: 499
}

func ExampleTimeMilli_Round() {
	t := xtime.DateMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC)
	round := []time.Duration{
		2 * time.Millisecond,
		time.Second,
		2 * time.Second,
		time.Minute,
		10 * time.Minute,
		time.Hour,
	}

	for _, d := range round {
		fmt.Printf("t.Round(%6s) = %s\n", d, t.Round(d).Format("15:04:05.999"))
	}
	// Output:
	// t.Round(   2ms) = 21:12:01.5
	// t.Round(    1s) = 21:12:01
	// t.Round(    2s) = 21:12:02
	// t.Round(  1m0s) = 21:12:00
	// t.Round( 10m0s) = 21:10:00
	// t.Round(1h0m0s) = 21:00:00
}

func ExampleTimeMilli_Truncate() {
	t := xtime.DateMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC)
	trunc := []time.Duration{
		2 * time.Millisecond,
		time.Second,
		2 * time.Second,
		time.Minute,
		10 * time.Minute,
	}

	for _, d := range trunc {
		fmt.Printf("t.Truncate(%5s) = %s\n", d, t.Truncate(d).Format("15:04:05.999"))
	}

	// Output:
	// t.Truncate(  2ms) = 21:12:01.498
	// t.Truncate(   1s) = 21:12:01
	// t.Truncate(   2s) = 21:12:00
	// t.Truncate( 1m0s) = 21:12:00
	// t.Truncate(10m0s) = 21:10:00
}

func ExampleTimeMilli_Unix() {
	t := xtime.DateMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC)
	fmt.Printf("%d\n", t.UnixMilli()) // milliseconds since 1970
	// Output: 1468185121499
}

func ExampleDateStampMilli() {
	t := xtime.DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	fmt.Printf("Go launched at %s\n", t)
	// Output: Go launched at 2016-07-10 21:12:00.499 +0000 UTC
}

func ExampleParseStampMilli() {
	t, err := xtime.ParseStampMilli(xtime.RFC3339Milli, "2016-07-10T21:12:00.499Z")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", t)
	// Output: 2016-07-10 21:12:00.499 +0000 UTC
}

func ExampleParseStampMilliInLocation() {
	t, err := xtime.ParseStampMilliInLocation(xtime.RFC3339Milli, "2016-07-10T21:12:00.499+02:00", time.FixedZone("CET", 2*60*60))
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", t)
	// Output: 2016-07-10 21:12:00.499 +0200 CET
}

func ExampleToStampMilli() {
	t := xtime.ToStampMilli(time.Date(2016, time.July, 10, 21, 12, 0, 499000000, time.UTC))
	fmt.Printf("Go launched at %s\n", t)
	// Output: Go launched at 2016-07-10 21:12:00.499 +0000 UTC
}

func ExampleTimestampMilli_Add() {
	start := xtime.DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	afterTenSeconds := start.Add(time.Second * 10)
	afterTenMinutes := start.Add(time.Minute * 10)
	afterTenHours := start.Add(time.Hour * 10)
	afterTenDays := start.Add(time.Hour * 24 * 10)

	fmt.Printf("start = %v\n", start)
	fmt.Printf("start.Add(time.Second * 10) = %s\n", afterTenSeconds)
	fmt.Printf("start.Add(time.Minute * 10) = %s\n", afterTenMinutes)
	fmt.Printf("start.Add(time.Hour * 10) = %s\n", afterTenHours)
	fmt.Printf("start.Add(time.Hour * 24 * 10) = %s\n", afterTenDays)

	// Output:
	// start = 2016-07-10 21:12:00.499 +0000 UTC
	// start.Add(time.Second * 10) = 2016-07-10 21:12:10.499 +0000 UTC
	// start.Add(time.Minute * 10) = 2016-07-10 21:22:00.499 +0000 UTC
	// start.Add(time.Hour * 10) = 2016-07-11 07:12:00.499 +0000 UTC
	// start.Add(time.Hour * 24 * 10) = 2016-07-20 21:12:00.499 +0000 UTC
}

func ExampleTimestampMilli_AddDate() {
	start := xtime.DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	oneDayLater := start.AddDate(0, 0, 1)
	oneMonthLater := start.AddDate(0, 1, 0)
	oneYearLater := start.AddDate(1, 0, 0)

	fmt.Printf("oneDayLater: start.AddDate(0, 0, 1) = %s\n", oneDayLater)
	fmt.Printf("oneMonthLater: start.AddDate(0, 1, 0) = %s\n", oneMonthLater)
	fmt.Printf("oneYearLater: start.AddDate(1, 0, 0) = %s\n", oneYearLater)

	// Output:
	// oneDayLater: start.AddDate(0, 0, 1) = 2016-07-11 21:12:00.499 +0000 UTC
	// oneMonthLater: start.AddDate(0, 1, 0) = 2016-08-10 21:12:00.499 +0000 UTC
	// oneYearLater: start.AddDate(1, 0, 0) = 2017-07-10 21:12:00.499 +0000 UTC
}

func ExampleTimestampMilli_In() {
	t := xtime.DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	t = t.In(time.FixedZone("CET", 2*60*60))
	fmt.Printf("%s\n", t)
	// Output: 2016-07-10 23:12:00.499 +0200 CET
}

func ExampleTimestampMilli_Millisecond() {
	t := xtime.DateStampMilli(2016, time.July, 10, 21, 12, 0, 499, time.UTC)
	msec := t.Millisecond()
	fmt.Printf("msec: %d\n", msec)
	// Output: msec: 499
}

func ExampleTimestampMilli_Round() {
	t := xtime.DateStampMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC)
	round := []time.Duration{
		2 * time.Millisecond,
		time.Second,
		2 * time.Second,
		time.Minute,
		10 * time.Minute,
		time.Hour,
	}

	for _, d := range round {
		fmt.Printf("t.Round(%6s) = %s\n", d, t.Round(d).Format("15:04:05.999"))
	}
	// Output:
	// t.Round(   2ms) = 21:12:01.5
	// t.Round(    1s) = 21:12:01
	// t.Round(    2s) = 21:12:02
	// t.Round(  1m0s) = 21:12:00
	// t.Round( 10m0s) = 21:10:00
	// t.Round(1h0m0s) = 21:00:00
}

func ExampleTimestampMilli_Truncate() {
	t := xtime.DateStampMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC)
	trunc := []time.Duration{
		2 * time.Millisecond,
		time.Second,
		2 * time.Second,
		time.Minute,
		10 * time.Minute,
	}

	for _, d := range trunc {
		fmt.Printf("t.Truncate(%5s) = %s\n", d, t.Truncate(d).Format("15:04:05.999"))
	}

	// Output:
	// t.Truncate(  2ms) = 21:12:01.498
	// t.Truncate(   1s) = 21:12:01
	// t.Truncate(   2s) = 21:12:00
	// t.Truncate( 1m0s) = 21:12:00
	// t.Truncate(10m0s) = 21:10:00
}

func ExampleTimestampMilli_Unix() {
	t := xtime.DateStampMilli(2016, time.July, 10, 21, 12, 1, 499, time.UTC)
	fmt.Printf("%d\n", t.UnixMilli()) // milliseconds since 1970
	// Output: 1468185121499
}
