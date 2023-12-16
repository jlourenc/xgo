// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xhttptrace_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jlourenc/xgo/xnet/xhttp/xhttptrace"
)

func Example() {
	ctx := context.Background()
	ctx = xhttptrace.WithClientTrace(ctx, &xhttptrace.ClientTrace{
		Retry: func(ri xhttptrace.RetryInfo) {
			fmt.Printf("retry count: %d, status code: %d", ri.RetryCount, ri.StatusCode)
		},
	})

	trace := xhttptrace.ContextClientTrace(ctx)
	trace.Retry(xhttptrace.RetryInfo{
		RetryCount: 1,
		StatusCode: http.StatusTooManyRequests,
	})

	// Output: retry count: 1, status code: 429
}
