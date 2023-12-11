// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xhttptrace

import (
	"context"
)

type (
	// ClientTrace is a set of hooks to run at various stages of an outgoing
	// HTTP request. Any particular hook may be nil.
	ClientTrace struct {
		// Retry is called before a round trip retry is made.
		Retry func(RetryInfo)
	}

	// RetryInfo contains information about the HTTP request retry.
	RetryInfo struct {
		// RetryCount is the retry count for a given HTTP request.
		RetryCount int

		// StatusCode specifies the HTTP response code gotten before trigering a retry.
		StatusCode int
	}

	clientEventContextKey struct{}
)

// ContextClientTrace returns the ClientTrace value stored in ctx.
// If none, it returns nil.
func ContextClientTrace(ctx context.Context) *ClientTrace {
	trace, _ := ctx.Value(clientEventContextKey{}).(*ClientTrace) //nolint:errcheck,revive // nil returned if none.
	return trace
}

// WithClientTrace returns a new context based on the provided parent ctx.
// HTTP client requests made with the returned context will use
// the provided trace hooks.
func WithClientTrace(parent context.Context, trace *ClientTrace) context.Context {
	if trace == nil {
		return parent
	}
	return context.WithValue(parent, clientEventContextKey{}, trace)
}
