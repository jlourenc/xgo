// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xhttp

import "net/http"

func isRequestIdempotent(req *http.Request) bool {
	switch req.Method {
	// idempotent methods: https://datatracker.ietf.org/doc/html/rfc9110#section-9.2.2
	case http.MethodDelete, http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPut, http.MethodTrace:
		return true
	default:
		// Any request is retryable if either Idempotency-Key or X-Idempotency-Key request header is set.
		return req.Header.Get(HeaderIdempotencyKey) != "" || req.Header.Get(HeaderXIdempotencyKey) != ""
	}
}

func isRequestRewindable(req *http.Request) bool {
	return req.Body == nil || req.Body == http.NoBody || req.GetBody != nil
}
