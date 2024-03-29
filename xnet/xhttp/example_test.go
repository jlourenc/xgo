// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xhttp_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jlourenc/xgo/xnet/xhttp"
	"github.com/jlourenc/xgo/xnet/xhttp/xhttptrace"
)

func ExampleHeaderExist() {
	headers := http.Header{
		"Header-Key": {"key1=val1", "key2", "key3=val3, key4"},
	}

	exist := xhttp.HeaderExist(headers, "Header-Key")

	fmt.Printf("got: %t", exist)
	// Output: got: true
}

func ExampleHeaderKeyValues() {
	headers := http.Header{
		"Header-Key": {"key1=val1", "key2", "key3=val3, key4"},
	}

	headerKeyValues := xhttp.HeaderKeyValues(headers, "Header-Key")

	fmt.Printf("got: %s", headerKeyValues)
	// Output: got: map[key1:val1 key2: key3:val3 key4:]
}

func ExampleHeaderValues() {
	headers := http.Header{
		"Header-Key": {"key1=val1", "key2", "key3=val3, key4"},
	}

	headerValues := xhttp.HeaderValues(headers, "Header-Key")

	fmt.Printf("got: %s", headerValues)
	// Output: got: [key1=val1 key2 key3=val3 key4]
}

func ExampleNewRetryTransport() {
	client := http.Client{
		Transport: xhttp.NewRetryTransport(
			xhttp.RetryTransportInitialInterval(100*time.Millisecond),
			xhttp.RetryTransportIntervalMultiplier(1.2),
			xhttp.RetryTransportJitterFactor(0.1),
			xhttp.RetryTransportMaxInterval(10*time.Second),
		),
		Timeout: 30 * time.Second,
	}

	ctx := xhttptrace.WithClientTrace(context.Background(), &xhttptrace.ClientTrace{
		Retry: func(ri xhttptrace.RetryInfo) {
			fmt.Printf("retry count: %d, status code: %d", ri.RetryCount, ri.StatusCode)
		},
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com", http.NoBody)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = client.Do(req); err != nil {
		log.Fatal(err)
	}
}

func ExampleParseHeaderDate() {
	headers := http.Header{
		xhttp.HeaderDate: []string{"Sun, 10 Jul 2016 21:12:00.499 GMT"},
	}

	d, err := xhttp.ParseHeaderDate(headers)
	if err != nil {
		log.Fatalf("an unexpected error occurred: %s", err)
	}

	fmt.Printf("date: %s", d)
	// Output: date: 2016-07-10 21:12:00.499 +0000 UTC
}

func ExampleReplaceHeader() {
	headers := http.Header{
		"Header-Key": {"key1=val1", "key2", "key3=val3, key4"},
	}

	xhttp.ReplaceHeader(headers, "prefix", "header-key", "key5=val5", "key6")
	fmt.Printf("got: %s\n", headers)

	xhttp.ReplaceHeader(headers, "prefix", "header-key", "key7=val7", "key8")
	fmt.Printf("got: %s\n", headers)

	// Output:
	// got: map[Header-Key:[key5=val5 key6] Prefix-Header-Key:[key1=val1 key2 key3=val3, key4]]
	// got: map[Header-Key:[key7=val7 key8] Prefix-1-Header-Key:[key1=val1 key2 key3=val3, key4] Prefix-Header-Key:[key5=val5 key6]]
}
