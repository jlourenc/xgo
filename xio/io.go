// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xio extends the Go standard library package io
// by providing additional I/O primitives.
package xio

import (
	"bytes"
	"io"
)

// DrainClose discards the entire ReadCloser and closes it.
func DrainClose(rc io.ReadCloser) error {
	if rc == nil {
		return nil
	}

	_, derr := io.Copy(io.Discard, rc)

	if cerr := rc.Close(); cerr != nil {
		return cerr
	}

	return derr
}

// DuplicateReader reads all of b to memmory and then returns two equivalent
// Readers yielding the same bytes.
//
// It returns an error if the initial slurp of all bytes fails. It does not attempt
// to make the returned ReadClosers have identical error-matching behavior.
func DuplicateReader(r io.Reader) (r1, r2 io.Reader, err error) {
	if r == nil {
		return nil, nil, nil
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(r); err != nil {
		return nil, r, err
	}

	return &buf, bytes.NewReader(buf.Bytes()), nil
}

// DuplicateReadCloser reads all of b to memory and then returns two equivalent
// ReadClosers yielding the same bytes.
//
// It returns an error if the initial slurp of all bytes fails. It does not attempt
// to make the returned ReadClosers have identical error-matching behavior.
func DuplicateReadCloser(rc io.ReadCloser) (rc1, rc2 io.ReadCloser, err error) {
	if rc == nil {
		return nil, nil, nil
	}

	var buf bytes.Buffer
	if _, err = buf.ReadFrom(rc); err != nil {
		return nil, rc, err
	}

	if err := rc.Close(); err != nil {
		return nil, rc, err
	}

	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
