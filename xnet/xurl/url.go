// Copyright 2022 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xurl extends the Go standard library package url
// by providing additional primitives and structures.
package xurl

import (
	"net/url"
	"path"
	"strings"
)

const (
	slash = "/"
)

// JoinBasePath joins a base path and any number of path elements into a single path, escaping path elements and
// separating them as well as the base with slashes. Empty elements are ignored. If the argument list is empty or
// all its elements are empty, JoinBasePath returns the base path only.
func JoinBasePath(base string, elems ...string) string {
	if !strings.HasSuffix(base, slash) {
		base += slash
	}

	return base + JoinPath(elems...)
}

// JoinPath joins any number of path elements into a single path, escaping and separating them with slashes.
// Empty elements are ignored. If the argument list is empty or all its elements are empty,
// JoinPath returns an empty string.
func JoinPath(elems ...string) string {
	escapedElems := make([]string, len(elems))
	for i, e := range elems {
		escapedElems[i] = url.PathEscape(e)
	}

	return path.Join(escapedElems...)
}
