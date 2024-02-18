// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xhttp_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jlourenc/xgo/xnet/xhttp"
	"github.com/jlourenc/xgo/xnet/xhttp/xhttptrace"
)

func TestRetryTransport_RoundTrip(t *testing.T) {
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	traceRetryCalled := false
	traceCtx := xhttptrace.WithClientTrace(context.Background(), &xhttptrace.ClientTrace{
		Retry: func(ri xhttptrace.RetryInfo) {
			traceRetryCalled = true
			if ri.RetryCount <= 0 {
				t.Errorf("retry count less than or equal to 0: %d", ri.RetryCount)
			}
			if ri.StatusCode != http.StatusServiceUnavailable {
				t.Errorf("status code mismatch: expected %d; got %d", http.StatusServiceUnavailable, ri.StatusCode)
			}
		},
	})
	defer func() {
		if !traceRetryCalled {
			t.Error("trace retry function did not get called")
		}
	}()

	u, _ := url.Parse("http://example.com")
	resp204 := &http.Response{StatusCode: http.StatusNoContent}
	resp413 := &http.Response{
		Header:     http.Header{xhttp.HeaderRetryAfter: []string{"1"}},
		StatusCode: http.StatusRequestEntityTooLarge,
	}
	resp429 := &http.Response{
		Header:     http.Header{xhttp.HeaderRetryAfter: []string{time.Now().Add(50 * time.Millisecond).Format(http.TimeFormat)}},
		StatusCode: http.StatusTooManyRequests,
	}
	resp503 := &http.Response{StatusCode: http.StatusServiceUnavailable}

	testCases := []struct {
		name         string
		ctx          context.Context //nolint:containedctx // ctx appended to req object
		jitterFactor float64
		next         http.RoundTripper
		req          *http.Request
		expectedResp *http.Response
		expectedErr  error
	}{
		{
			name: "next round trip error bubbles up as is",
			next: &fakeTransport{},
			req: &http.Request{
				Body:   http.NoBody,
				Method: http.MethodGet,
				URL:    u,
			},
			expectedErr: errNoResponse,
		},
		{
			name: "success response with no retry",
			next: &fakeTransport{resps: []*http.Response{resp204}},
			req: &http.Request{
				Body:   io.NopCloser(strings.NewReader("payload")),
				Method: http.MethodPut,
				URL:    u,
			},
			expectedResp: resp204,
		},
		{
			name: "success response with a single retry on 413 with retry-after header",
			next: &fakeTransport{resps: []*http.Response{resp413, resp204}},
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("payload")),
				GetBody: func() (io.ReadCloser, error) {
					return io.NopCloser(strings.NewReader("payload")), nil
				},
				Header: http.Header{
					xhttp.HeaderIdempotencyKey: []string{"idempotency-key-id"},
				},
				Method: http.MethodPost,
				URL:    u,
			},
			expectedResp: resp204,
		},
		{
			name: "success response with a single retry on 429",
			next: &fakeTransport{resps: []*http.Response{resp429, resp204}},
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("payload")),
				GetBody: func() (io.ReadCloser, error) {
					return io.NopCloser(strings.NewReader("payload")), nil
				},
				Header: http.Header{
					xhttp.HeaderXIdempotencyKey: []string{"idempotency-key-id"},
				},
				Method: http.MethodPost,
				URL:    u,
			},
			expectedResp: resp204,
		},
		{
			name: "success response with multiple retries on 503",
			next: &fakeTransport{resps: []*http.Response{resp503, resp503, resp503, resp204}},
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("payload")),
				GetBody: func() (io.ReadCloser, error) {
					return io.NopCloser(strings.NewReader("payload")), nil
				},
				Method: http.MethodPut,
				URL:    u,
			},
			jitterFactor: 0.5,
			expectedResp: resp204,
		},
		{
			name: "failure to rewing body returns last http response",

			next: &fakeTransport{resps: []*http.Response{resp429, resp204}},
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("payload")),
				GetBody: func() (io.ReadCloser, error) {
					return nil, errNoBody
				},
				Method: http.MethodPut,
				URL:    u,
			},
			expectedResp: resp429,
		},
		{
			name: "context done cancels retry and returns last http response",
			ctx:  canceledCtx,
			next: &fakeTransport{resps: []*http.Response{resp429, resp204}},
			req: &http.Request{
				Body:   http.NoBody,
				Method: http.MethodDelete,
				URL:    u,
			},
			expectedResp: resp429,
		},
		{
			name: "trace retry called on retries",
			ctx:  traceCtx,
			next: &fakeTransport{resps: []*http.Response{resp503, resp503, resp204}},
			req: &http.Request{
				Body:   http.NoBody,
				Method: http.MethodDelete,
				URL:    u,
			},
			expectedResp: resp204,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			retryTransp := xhttp.NewRetryTransport(
				xhttp.RetryTransportNextRoundTripper(tc.next),
				xhttp.RetryTransportInitialInterval(10*time.Millisecond),
				xhttp.RetryTransportIntervalMultiplier(2),
				xhttp.RetryTransportJitterFactor(tc.jitterFactor),
				xhttp.RetryTransportMaxInterval(20*time.Millisecond))

			if tc.ctx != nil {
				tc.req = tc.req.Clone(tc.ctx)
			}

			gotResp, gotErr := retryTransp.RoundTrip(tc.req)

			if gotResp != tc.expectedResp {
				t.Errorf("response mistmatch: %v != %v", gotResp, tc.expectedResp)
			}
			if gotErr != tc.expectedErr {
				t.Errorf("error mistmach: %v != %v", gotErr, tc.expectedErr)
			}
		})
	}
}

func TestRetryTransportInitialInterval(t *testing.T) {
	testCases := []struct {
		name     string
		interval time.Duration
		panic    bool
	}{
		{
			name:     "panic",
			interval: 0,
			panic:    true,
		},
		{
			name:     "valid",
			interval: 1,
			panic:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testRetryTransportOptionPanic(t, tc.panic, func() xhttp.RetryTransportOption {
				return xhttp.RetryTransportInitialInterval(tc.interval)
			})
		})
	}
}

func TestRetryTransportIntervalMultiplier(t *testing.T) {
	testCases := []struct {
		name       string
		multiplier float64
		panic      bool
	}{
		{
			name:       "panic",
			multiplier: 0.0,
			panic:      true,
		},
		{
			name:       "valid",
			multiplier: 1.0,
			panic:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testRetryTransportOptionPanic(t, tc.panic, func() xhttp.RetryTransportOption {
				return xhttp.RetryTransportIntervalMultiplier(tc.multiplier)
			})
		})
	}
}

func TestRetryTransportJitterFactor(t *testing.T) {
	testCases := []struct {
		name   string
		factor float64
		panic  bool
	}{
		{
			name:   "panic - below range",
			factor: -0.1,
			panic:  true,
		},
		{
			name:   "panic - above range",
			factor: 1.1,
			panic:  true,
		},
		{
			name:   "valid",
			factor: 0.5,
			panic:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testRetryTransportOptionPanic(t, tc.panic, func() xhttp.RetryTransportOption {
				return xhttp.RetryTransportJitterFactor(tc.factor)
			})
		})
	}
}

func TestRetryTransportMaxInterval(t *testing.T) {
	testCases := []struct {
		name     string
		interval time.Duration
		panic    bool
	}{
		{
			name:     "panic",
			interval: 0,
			panic:    true,
		},
		{
			name:     "valid",
			interval: 1,
			panic:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testRetryTransportOptionPanic(t, tc.panic, func() xhttp.RetryTransportOption {
				return xhttp.RetryTransportMaxInterval(tc.interval)
			})
		})
	}
}

func TestRetryTransportNextRoundTripper(t *testing.T) {
	testCases := []struct {
		name  string
		next  http.RoundTripper
		panic bool
	}{
		{
			name:  "panic",
			next:  nil,
			panic: true,
		},
		{
			name:  "valid",
			next:  &fakeTransport{},
			panic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testRetryTransportOptionPanic(t, tc.panic, func() xhttp.RetryTransportOption {
				return xhttp.RetryTransportNextRoundTripper(tc.next)
			})
		})
	}
}
