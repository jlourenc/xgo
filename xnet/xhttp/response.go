// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package xhttp

import "net/http"

func isResponseRetryable(resp *http.Response) bool {
	switch resp.StatusCode {
	// 4xx status codes
	case http.StatusRequestTimeout, http.StatusTooEarly, http.StatusTooManyRequests:
		return true
	// 413 is retryable if Retry-After header is set.
	case http.StatusRequestEntityTooLarge:
		return resp.Header.Get(HeaderRetryAfter) != ""
	// 5xx status codes
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
