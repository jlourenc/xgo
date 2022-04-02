// Copyright 2022 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xio_test

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jlourenc/xgo/xio"
)

func ExampleDrainClose() {
	rc := io.NopCloser(strings.NewReader("message"))

	if err := xio.DrainClose(rc); err != nil {
		log.Fatalf("Failed to drain and close ReadCloser: %v", err)
	}

	log.Print("ReadCloser drained and closed")
}

func ExampleDuplicateReader() {
	r := strings.NewReader("message")

	r1, r2, err := xio.DuplicateReader(r)
	if err != nil {
		log.Fatalf("Failed to duplicate Reader: %v", err)
	}

	b1, _ := io.ReadAll(r1)
	b2, _ := io.ReadAll(r2)

	fmt.Printf("%s\n", b1)
	fmt.Printf("%s\n", b2)
	// Output:
	// message
	// message
}

func ExampleDuplicateReadCloser() {
	rc := io.NopCloser(strings.NewReader("message"))

	rc1, rc2, err := xio.DuplicateReadCloser(rc)
	if err != nil {
		log.Fatalf("Failed to duplicate Reader: %v", err)
	}
	defer rc1.Close()
	defer rc2.Close()

	b1, _ := io.ReadAll(rc1)
	b2, _ := io.ReadAll(rc2)

	fmt.Printf("%s\n", b1)
	fmt.Printf("%s\n", b2)
	// Output:
	// message
	// message
}
