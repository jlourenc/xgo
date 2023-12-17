// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xhttp

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/jlourenc/xgo/xio"
	"github.com/jlourenc/xgo/xnet/xhttp/xhttptrace"
)

const (
	retryTransportDefaultInitialInterval    = 200 * time.Millisecond
	retryTransportDefaultIntervalMultiplier = 1.5
	retryTransportDefaultJitterFactor       = 0.2
	retryTransportDefaultMaxInterval        = 30 * time.Second
)

// RetryTransport is an HTTP transport that implements HTTP retries according to
// the HTTP semantics defined in https://datatracker.ietf.org/doc/html/rfc9110.
type retryTransport struct {
	next http.RoundTripper

	// backoff policy
	initialInterval    time.Duration
	intervalMultiplier float64
	jitterFactor       float64
	maxInterval        time.Duration
}

// NewRetryTransport creates a new RetryTransport configured with the options passed in input,
// notably the backoff policy and the next round tripper in the chain.
func NewRetryTransport(options ...RetryTransportOption) http.RoundTripper {
	t := &retryTransport{
		initialInterval:    retryTransportDefaultInitialInterval,
		intervalMultiplier: retryTransportDefaultIntervalMultiplier,
		jitterFactor:       retryTransportDefaultJitterFactor,
		maxInterval:        retryTransportDefaultMaxInterval,
		next:               http.DefaultTransport,
	}

	for _, opt := range options {
		opt.apply(t)
	}

	return t
}

// RoundTrip makes RetryTransport implement the RoundTripper interface.
//
// It retries retryable (as defined by their status code) responses of idempotent requests,
// following a backoff policy or respecting Retry-After response headers.
//
// See HTTP semantics defined in: https://datatracker.ietf.org/doc/html/rfc9110.
func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	reqRetryable := isRequestIdempotent(req) && isRequestRewindable(req)
	retryCount := 0
	retryInterval := t.initialInterval

	trace := xhttptrace.ContextClientTrace(ctx)
	if trace == nil {
		trace = &xhttptrace.ClientTrace{}
	}

	for {
		resp, err := t.next.RoundTrip(req)
		if err != nil {
			return resp, err
		}

		if !reqRetryable || !isResponseRetryable(resp) {
			return resp, nil
		}

		// Clone request if body is rewindable.
		if req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return resp, nil //nolint:nilerr // return last response
			}
			req = req.Clone(ctx)
			req.Body = body
		}

		timer := time.NewTimer(computeWaitDuration(retryInterval, t.jitterFactor, resp.Header))
		select {
		case <-ctx.Done():
			timer.Stop()
			return resp, nil
		case <-timer.C:
			xio.DrainClose(resp.Body)
		}

		retryInterval = time.Duration(float64(retryInterval) * t.intervalMultiplier)
		if retryInterval > t.maxInterval {
			retryInterval = t.maxInterval
		}
		retryCount++

		if trace.Retry != nil {
			trace.Retry(xhttptrace.RetryInfo{
				RetryCount: retryCount,
				StatusCode: resp.StatusCode,
			})
		}
	}
}

func computeWaitDuration(interval time.Duration, jitterFactor float64, headers http.Header) time.Duration {
	if retryAfter := headers.Get(HeaderRetryAfter); retryAfter != "" {
		if secs, err := strconv.Atoi(retryAfter); err == nil {
			return time.Duration(secs) * time.Second
		}

		if date, err := http.ParseTime(retryAfter); err == nil {
			return time.Until(date)
		}
	}

	if jitterFactor == 0.0 {
		return interval
	}

	delta := jitterFactor * float64(interval)
	minInterval := float64(interval) - delta

	// returns a random value in the half-open interval [interval - delta, interval + delta).
	return time.Duration(minInterval + (rand.Float64() * delta * 2)) //nolint:gosec // rand is used in a non security-sensitive scenario
}

type (
	// RetryTransportOption configures the RetryTransport options
	// when calling NewRetryTransport.
	RetryTransportOption interface {
		apply(d *retryTransport)
	}

	funcRetryTransportOption struct {
		fn func(*retryTransport)
	}
)

func newFuncRetryTransportOption(fn func(*retryTransport)) funcRetryTransportOption {
	return funcRetryTransportOption{
		fn: fn,
	}
}

func (o funcRetryTransportOption) apply(d *retryTransport) {
	o.fn(d)
}

// RetryTransportInitialInterval returns a RetryTransportOption that configures the
// initial retry interval of the backoff policy. Value must be > 0, otherwise it panics.
func RetryTransportInitialInterval(interval time.Duration) RetryTransportOption {
	if interval <= 0 {
		panic("invalid initial interval value")
	}
	return newFuncRetryTransportOption(func(rt *retryTransport) {
		rt.initialInterval = interval
	})
}

// RetryTransportIntervalMultiplier returns a RetryTransportOption that configures the
// interval multiplier of the backoff policy. Value must be >= 1.0, otherwise it panics.
func RetryTransportIntervalMultiplier(multiplier float64) RetryTransportOption {
	if multiplier < 1.0 {
		panic("invalid interval multiplier value")
	}
	return newFuncRetryTransportOption(func(rt *retryTransport) {
		rt.intervalMultiplier = multiplier
	})
}

// RetryTransportJitterFactor returns a RetryTransportOption that configures the jitter
// multiplier of the backoff policy to randomize the retry distribution. Value must be in the
// [0.0, 1.0] range, otherwise it panics.
func RetryTransportJitterFactor(factor float64) RetryTransportOption {
	if factor < 0.0 || factor > 1.0 {
		panic("invalid jitter factor value")
	}
	return newFuncRetryTransportOption(func(rt *retryTransport) {
		rt.jitterFactor = factor
	})
}

// RetryTransportMaxInterval returns a RetryTransportOption that configures the max interval of the
// backoff policy. Once reached, retry interval is not increased. Value must be > 0, otherwise it panics.
func RetryTransportMaxInterval(interval time.Duration) RetryTransportOption {
	if interval <= 0 {
		panic("invalid max interval value")
	}
	return newFuncRetryTransportOption(func(rt *retryTransport) {
		rt.maxInterval = interval
	})
}

// RetryTransportNextRoundTripper returns a RetryTransportOption that configures the
// next round tripper to call. If not used http.DefaultTransport will be used.
func RetryTransportNextRoundTripper(next http.RoundTripper) RetryTransportOption {
	if next == nil {
		panic("next http.RoundTripper is nil")
	}
	return newFuncRetryTransportOption(func(rt *retryTransport) {
		rt.next = next
	})
}
