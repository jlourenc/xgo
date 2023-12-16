// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xhttp extends the Go standard library package http
// by providing additional HTTP client and server implementations.
package xhttp

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTP standard headers.
const (
	// https://datatracker.ietf.org/doc/html/rfc7231#section-5.3.2
	HeaderAccept = "Accept"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-5.3.3
	HeacerAcceptCharset = "Accept-Charset"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-5.3.4
	HeaderAcceptEncoding = "Accept-Encoding"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-5.3.5
	HeaderAcceptLanguage = "Accept-Language"
	// https://datatracker.ietf.org/doc/html/rfc5789#section-3.1
	HeaderAcceptPath = "Accept-Patch"
	// https://www.w3.org/TR/ldp/#header-accept-post
	HeaderAcceptPost = "Accept-Post"
	// https://datatracker.ietf.org/doc/html/rfc7233#section-2.3
	HeaderAcceptRanges = "Accept-Ranges"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
	HeaderAccessControlExposeHeaders = "Access-Control-Expose-Headers"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	HeaderAccessControlMaxAge = "Access-Control-Max-Age"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method
	HeaderAccessControlRequestMethod = "Access-Control-Request-Method"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.1
	HeaderAge = "Age"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-7.4.1
	HeaderAllow = "Allow"
	// https://datatracker.ietf.org/doc/html/rfc7838#section-3
	HeaderAltSvc = "Alt-Svc"
	// https://datatracker.ietf.org/doc/html/rfc7235#section-4.2
	HeaderAuthorization = "Authorization"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2
	HeaderCacheControl = "Cache-Control"
	// https://www.w3.org/TR/clear-site-data/#header
	HeaderClearSiteData = "Clear-Site-Data"
	// https://datatracker.ietf.org/doc/html/rfc7230#section-6.1
	HeaderConnection = "Connection"
	// https://datatracker.ietf.org/doc/html/rfc6266#section-4
	HeaderContentDisposition = "Content-Disposition"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-3.1.2.2
	HeaderContentEncoding = "Content-Encoding"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-3.1.2.2
	HeaderContentLanguage = "Content-Language"
	// https://datatracker.ietf.org/doc/html/rfc7230#section-3.3.2
	HeaderContentLength = "Content-Length"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-3.1.4.2
	HeaderContentLocation = "Content-Location"
	// https://datatracker.ietf.org/doc/html/rfc7233#section-4.2
	HeaderContentRange = "Content-Range"
	// https://www.w3.org/TR/CSP3/#csp-header
	HeaderContentSecurityPolicy = "Content-Security-Policy"
	// https://www.w3.org/TR/CSP3/#cspro-header
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-3.1.1.5
	HeaderContentType = "Content-Type"
	// https://datatracker.ietf.org/doc/html/rfc6265#section-5.4
	HeaderCookie = "Cookie"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Embedder-Policy
	HeaderCrossOriginEmbedderPolicy = "Cross-Origin-Embedder-Policy"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Opener-Policy
	HeaderCrossOriginOpenerPolicy = "Cross-Origin-Opener-Policy"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cross-Origin-Resource-Policy
	HeaderCrossOriginResourcePolicy = "Cross-Origin-Resource-Policy"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-7.1.1.2
	HeaderDate = "Date"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Digest
	HeaderDigest = "Digest"
	// https://datatracker.ietf.org/doc/html/rfc7232#section-2.3
	HeaderEtag = "Etag"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-5.1.1
	HeaderExpect = "Expect"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Expect-CT
	HeaderExpectCT = "Expect-CT"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.3
	HeaderExpires = "Expires"
	// https://datatracker.ietf.org/doc/html/rfc7239#section-4
	HeaderForwarded = "Forwarded"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-5.5.1
	HeaderFrom = "From"
	// https://datatracker.ietf.org/doc/html/rfc7230#section-5.4
	HeaderHost = "Host"
	// https://datatracker.ietf.org/doc/html/rfc7540#section-3.2.1
	HeaderHTTP2Settings = "HTTP2-Settings"
	// https://datatracker.ietf.org/doc/html/rfc7232#section-3.1
	HeaderIfMatch = "If-Match"
	// https://datatracker.ietf.org/doc/html/rfc7232#section-3.3
	HeaderIfModifiedSince = "If-Modified-Since"
	// https://datatracker.ietf.org/doc/html/rfc7232#section-3.2
	HeaderIfNoneMatch = "If-None-Match"
	// https://datatracker.ietf.org/doc/html/rfc7233#section-3.2
	HeaderIfRange = "If-Range"
	// https://datatracker.ietf.org/doc/html/rfc7232#section-3.4
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Keep-Alive
	HeaderKeepAlive = "Keep-Alive"
	// https://datatracker.ietf.org/doc/html/rfc7232#section-2.2
	HeaderLastModified = "Last-Modified"
	// https://datatracker.ietf.org/doc/html/rfc5988#section-5
	HeaderLink = "Link"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-7.1.2
	HeaderLocation = "Location"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-5.1.2
	HeaderMaxForwards = "Max-Forwards"
	// https://www.w3.org/TR/network-error-logging/#nel-response-header
	HeaderNEL = "NEL"
	// https://datatracker.ietf.org/doc/html/rfc6454#section-7
	HeaderOrigin = "Origin"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.4
	HeaderPragma = "Pragma"
	// https://datatracker.ietf.org/doc/html/rfc7240#section-2
	HeaderPrefer = "Prefer"
	// https://datatracker.ietf.org/doc/html/rfc7235#section-4.3
	HeaderProxyAuthenticate = "Proxy-Authenticate"
	// https://datatracker.ietf.org/doc/html/rfc7235#section-4.4
	HeaderProxyAuthorization = "Proxy-Authorization"
	// https://datatracker.ietf.org/doc/html/rfc7233#section-3.1
	HeaderRange = "Range"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-5.5.2
	HeaderReferer = "Referer"
	// https://www.w3.org/TR/referrer-policy/
	HeaderReferrerPolicy = "Referrer-Policy"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-7.1.3
	HeaderRetryAfter = "Retry-After"
	// https://wicg.github.io/savedata/#save-data-request-header-field
	HeaderSaveData = "Save-Data"
	// https://www.w3.org/TR/fetch-metadata/#sec-fetch-dest-header
	HeaderSecFetchDest = "Sec-Fetch-Dest"
	// https://www.w3.org/TR/fetch-metadata/#sec-fetch-mode-header
	HeaderSecFetchMode = "Sec-Fetch-Mode"
	// https://www.w3.org/TR/fetch-metadata/#sec-fetch-site-header
	HeaderSecFetchSite = "Sec-Fetch-Site"
	// https://www.w3.org/TR/fetch-metadata/#sec-fetch-user-header
	HeaderSecFetchUser = "Sec-Fetch-User"
	// https://datatracker.ietf.org/doc/html/rfc6455#section-11.3.3
	HeaderSecWebSocketAccept = "Sec-WebSocket-Accept"
	// https://datatracker.ietf.org/doc/html/rfc6455#section-11.3.4
	HeaderSecWebSocketProtocol = "Sec-WebSocket-Protocol"
	// https://datatracker.ietf.org/doc/html/rfc6455#section-11.3.5
	HeaderSecWebSocketVersion = "Sec-WebSocket-Version"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-7.4.2
	HeaderServer = "Server"
	// https://www.w3.org/TR/server-timing/
	HeaderServerTiming = "Server-Timing"
	// https://datatracker.ietf.org/doc/html/rfc6265#section-5.2
	HeaderSetCookie = "Set-Cookie"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/SourceMap
	HeaderSourceMap = "SourceMap"
	// https://datatracker.ietf.org/doc/html/rfc6797#section-6.1
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	// https://datatracker.ietf.org/doc/html/rfc7230#section-4.3
	HeaderTE = "TE"
	// https://www.w3.org/TR/resource-timing-2/#sec-timing-allow-origin
	HeaderTimingAllowOrigin = "Timing-Allow-Origin"
	// https://datatracker.ietf.org/doc/html/rfc7230#section-4.4
	HeaderTrailer = "Trailer"
	// https://datatracker.ietf.org/doc/html/rfc7230#section-3.3.1
	HeaderTransferEncoding = "Transfer-Encoding"
	// https://datatracker.ietf.org/doc/html/rfc7230#section-6.7
	HeaderUpgrade = "Upgrade"
	// https://www.w3.org/TR/upgrade-insecure-requests/#preference
	HeaderUpgradeInsecureRequests = "Upgrade-Insecure-Requests"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-5.5.3
	HeaderUserAgent = "User-Agent"
	// https://datatracker.ietf.org/doc/html/rfc7231#section-7.1.4
	HeaderVary = "Vary"
	// https://datatracker.ietf.org/doc/html/rfc7230#section-5.7.1
	HeaderVia = "Via"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Want-Digest
	HeaderWantDigest = "Want-Digest"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.5
	HeaderWarning = "Warning"
	// https://datatracker.ietf.org/doc/html/rfc7235#section-4.1
	HeaderWWWAuthenticate = "WWW-Authenticate"
)

// HTTP non-standard headers, but widely used.
const (
	// https://datatracker.ietf.org/doc/draft-ietf-httpapi-idempotency-key-header/
	HeaderIdempotencyKey = "Idempotency-Key"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options
	HeaderXContentTypeOptions = "X-Content-Type-Options"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-DNS-Prefetch-Control
	HeaderXDNSPrefetchControl = "X-DNS-Prefetch-Control"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For
	HeaderXForwardedFor = "X-Forwarded-For"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-Host
	HeaderXForwardedHost = "X-Forwarded-Host"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-Proto
	HeaderXForwardedProto = "X-Forwarded-Proto"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
	HeaderXFrameOptions = "X-Frame-Options"
	// https://tools.ietf.org/id/draft-idempotency-header-01.html
	// Deprecated: use HeaderIdempotencyKey instead.
	HeaderXIdempotencyKey = "X-Idempotency-Key"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
	HeaderXXSSProtection = "X-XSS-Protection"
)

// HTTP Cache-Control directives.
const (
	// https://datatracker.ietf.org/doc/html/rfc8246#section-2
	CacheControlImmutable = "immutable"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.1.1 & https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2.8
	CacheControlMaxAge = "max-age"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.1.2
	CacheControlMaxStale = "max-stale"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.1.3
	CacheControlMinFresh = "min-fresh"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2.1
	CacheControlMustRevalidate = "must-revalidate"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.1.4 & https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2.2
	CacheControlNoCache = "no-cache"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.1.5 & https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2.3
	CacheControlNoStore = "no-store"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.1.6 & https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2.4
	CacheControlNoTransform = "no-transform"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.1.7
	CacheControlOnlyIfCached = "only-if-cached"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2.5
	CacheControlPublic = "public"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2.6
	CacheControlPrivate = "private"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2.7
	CacheControlProxyRevalidate = "proxy-revalidate"
	// https://datatracker.ietf.org/doc/html/rfc7234#section-5.2.2.9
	CacheControlSMaxAge = "s-maxage"
	// https://datatracker.ietf.org/doc/html/rfc5861#section-4
	CacheControlStaleIfError = "stale-if-error"
	// https://datatracker.ietf.org/doc/html/rfc5861#section-3
	CacheControlStaleWhileRevalidate = "stale-while-revalidate"
)

var errHeaderNoDate = errors.New("no date header")

// HeaderExist returns whether the key exists in headers.
func HeaderExist(headers http.Header, key string) bool {
	_, ok := headers[http.CanonicalHeaderKey(key)]
	return ok
}

// HeaderKeyValues returns all key/value pairs associated with the given key, or nil if the key does not exist.
// It is case insensitive; textproto.CanonicalMIMEHeaderKey is used to canonicalize the provided key.
func HeaderKeyValues(headers http.Header, key string) map[string]string {
	values := HeaderValues(headers, key)
	if values == nil {
		return nil
	}

	headerKeyValues := make(map[string]string, len(values))
	for _, value := range values {
		if keyValue := strings.Split(value, "="); len(keyValue) > 1 {
			headerKeyValues[keyValue[0]] = keyValue[1]
		} else {
			headerKeyValues[value] = ""
		}
	}

	return headerKeyValues
}

// HeaderValues returns all values associated with the given key, or nil if the key does not exist.
// It is case insensitive; textproto.CanonicalMIMEHeaderKey is used to canonicalize the provided key.
// As per Section 4.2 of the RFC 2616 (http://www.w3.org/Protocols/rfc2616/rfc2616-sec4.html#sec4.2),
// values from multiple occurrences of a header should be concatenated, if the header's value is a comma-separated list.
func HeaderValues(headers http.Header, key string) (headerValues []string) {
	for _, value := range headers.Values(key) {
		fields := strings.Split(value, ",")
		for i, f := range fields {
			fields[i] = strings.TrimSpace(f)
		}
		headerValues = append(headerValues, fields...)
	}
	return headerValues
}

// ParseHeaderDate parses the Date header and returns its value as a time.Time if valid.
// An error is returned otherwise.
// https://datatracker.ietf.org/doc/html/rfc7231#section-7.1.1.1
func ParseHeaderDate(headers http.Header) (time.Time, error) {
	date := headers.Get(HeaderDate)
	if date == "" {
		return time.Time{}, errHeaderNoDate
	}
	return http.ParseTime(date)
}

// ReplaceHeader sets the values for the key in the headers. If the key already exists, the old values
// are preserved in a new key prefixed with prefix + either '-' or '-#-' (with # a strictly positive integer)
// depending on the exitence of the prefixed keys.
// https://www.w3.org/TR/ct-guidelines/#sec-original-headers
func ReplaceHeader(headers http.Header, prefix, key string, values ...string) {
	if headers == nil {
		return
	}

	prefix = http.CanonicalHeaderKey(prefix)
	key = http.CanonicalHeaderKey(key)
	prefixedKey := prefix + "-" + key

	if v, ok := headers[prefixedKey]; ok {
		i := 1
		for {
			tmp, ok := headers[prefix+"-"+strconv.Itoa(i)+"-"+key]
			headers[prefix+"-"+strconv.Itoa(i)+"-"+key] = v
			if !ok {
				break
			}
			v = tmp
			i++
		}
	}

	if v, ok := headers[key]; ok {
		headers[prefixedKey] = v
	}

	headers[http.CanonicalHeaderKey(key)] = values
}
