// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xurl_test

import (
	"fmt"

	"github.com/jlourenc/xgo/xnet/xurl"
)

func ExampleJoinBasePath() {
	path := xurl.JoinBasePath("https://pkg.go.dev/", "github.com", "jlourenc", "xgo")

	fmt.Printf("%s\n", path)
	// Output: https://pkg.go.dev/github.com/jlourenc/xgo
}

func ExampleJoinPath() {
	path := xurl.JoinPath("github.com", "jlourenc", "xgo")

	fmt.Printf("%s\n", path)
	// Output: github.com/jlourenc/xgo
}
