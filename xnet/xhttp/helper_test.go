// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xhttp_test

import (
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/jlourenc/xgo/xnet/xhttp"
)

var (
	errNoBody     = errors.New("no body")
	errNoResponse = errors.New("no response")
)

type fakeTransport struct {
	counter   int
	reqBodies [][]byte
	resps     []*http.Response
}

func (t *fakeTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	var b []byte
	if req.Body == nil {
		return nil, errNoBody
	}

	b, err = io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	t.reqBodies = append(t.reqBodies, b)

	if t.counter >= len(t.resps) {
		return nil, errNoResponse
	}

	resp = t.resps[t.counter]
	t.counter++

	return resp, nil
}

func testRetryTransportOptionPanic(tb testing.TB, shouldPanic bool, fn func() xhttp.RetryTransportOption) {
	tb.Helper()

	defer func() {
		r := recover()
		if r == nil && shouldPanic {
			tb.Error("panic expected; got none")
		} else if r != nil && !shouldPanic {
			tb.Error("no panic expected; got one")
		}
	}()

	got := fn()

	if got == nil {
		tb.Error("retry option expected; got nil")
	}
}
