// Copyright 2021 Jérémy Lourenço. All rights reserved.
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
