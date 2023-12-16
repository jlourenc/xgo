// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xhttptrace_test

import (
	"context"
	"testing"

	"github.com/jlourenc/xgo/xnet/xhttp/xhttptrace"
)

func TestContextClientTrace(t *testing.T) {
	testCases := []struct {
		name          string
		trace         *xhttptrace.ClientTrace
		expectedTrace *xhttptrace.ClientTrace
		expectedOK    bool
	}{
		{
			name:  "nil trace",
			trace: nil,
		},
		{
			name: "specified trace",
			trace: &xhttptrace.ClientTrace{
				Retry: func(ri xhttptrace.RetryInfo) {},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := xhttptrace.WithClientTrace(context.Background(), tc.trace)

			got := xhttptrace.ContextClientTrace(ctx)

			if tc.trace != got {
				t.Errorf("expected %v; got %v", tc.trace, got)
			}
		})
	}
}
