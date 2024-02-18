// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xio_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/jlourenc/xgo/xio"
)

func TestDrainClose(t *testing.T) {
	testCases := []struct {
		name     string
		rc       io.ReadCloser
		expected error
	}{
		{
			name:     "nil ReadCloser",
			rc:       nil,
			expected: nil,
		},
		{
			name:     "failed to close",
			rc:       errClose{},
			expected: io.ErrUnexpectedEOF,
		},
		{
			name:     "failed to read",
			rc:       errRead{},
			expected: io.ErrNoProgress,
		},
		{
			name:     "failed to read and close",
			rc:       errReadClose{},
			expected: io.ErrUnexpectedEOF,
		},
		{
			name:     "succeeded to read and close",
			rc:       io.NopCloser(bytes.NewReader([]byte(`message`))),
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := xio.DrainClose(tc.rc)

			if tc.expected != got {
				t.Errorf("expected %s; got %s", tc.expected, got)
			}
		})
	}
}

func TestDuplicateReader(t *testing.T) {
	testCases := []struct {
		name        string
		r           io.Reader
		expectedB1  []byte
		expectedB2  []byte
		expectedErr error
	}{
		{
			name:        "nil Read",
			r:           nil,
			expectedB1:  nil,
			expectedB2:  nil,
			expectedErr: nil,
		},
		{
			name:        "failed to read",
			r:           errRead{},
			expectedB1:  nil,
			expectedB2:  []byte{},
			expectedErr: io.ErrNoProgress,
		},
		{
			name:        "succeeded to read",
			r:           bytes.NewReader([]byte(`message`)),
			expectedB1:  []byte(`message`),
			expectedB2:  []byte(`message`),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r1, r2, err := xio.DuplicateReader(tc.r)

			if r1 != nil {
				b1, _ := io.ReadAll(r1)
				if tc.expectedB1 == nil {
					t.Errorf("expected r1 to be nil; got %q", b1)
				} else if !bytes.Equal(b1, tc.expectedB1) {
					t.Errorf("expected r21 %q; got %q", tc.expectedB1, b1)
				}
			} else if tc.expectedB1 != nil {
				t.Errorf("expected r21 %q; got nil", tc.expectedB1)
			}

			if r2 != nil {
				b2, _ := io.ReadAll(r2)
				if tc.expectedB2 == nil {
					t.Errorf("expected r2 to be nil; got %q", b2)
				} else if !bytes.Equal(b2, tc.expectedB2) {
					t.Errorf("expected r2 %q; got %q", tc.expectedB1, b2)
				}
			} else if tc.expectedB2 != nil {
				t.Errorf("expected r2 %q; got nil", tc.expectedB2)
			}

			if tc.expectedErr != err {
				t.Errorf("expected %s; got %s", tc.expectedErr, err)
			}
		})
	}
}

func TestDuplicateReadCloser(t *testing.T) {
	testCases := []struct {
		name        string
		rc          io.ReadCloser
		expectedB1  []byte
		expectedB2  []byte
		expectedErr error
	}{
		{
			name:        "nil ReadCloser",
			rc:          nil,
			expectedB1:  nil,
			expectedB2:  nil,
			expectedErr: nil,
		},
		{
			name:        "failed to read",
			rc:          errRead{},
			expectedB1:  nil,
			expectedB2:  []byte{},
			expectedErr: io.ErrNoProgress,
		},
		{
			name:        "failed to close",
			rc:          errClose{},
			expectedB1:  nil,
			expectedB2:  []byte{},
			expectedErr: io.ErrUnexpectedEOF,
		},
		{
			name:        "succeeded to read and close",
			rc:          io.NopCloser(bytes.NewReader([]byte(`message`))),
			expectedB1:  []byte(`message`),
			expectedB2:  []byte(`message`),
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rc1, rc2, err := xio.DuplicateReadCloser(tc.rc)

			if rc1 != nil {
				b1, _ := io.ReadAll(rc1)
				if tc.expectedB1 == nil {
					t.Errorf("expected rc1 to be nil; got %q", b1)
				} else if !bytes.Equal(b1, tc.expectedB1) {
					t.Errorf("expected rc1 %q; got %q", tc.expectedB1, b1)
				}
			} else if tc.expectedB1 != nil {
				t.Errorf("expected rc1 %q; got nil", tc.expectedB1)
			}

			if rc2 != nil {
				b2, _ := io.ReadAll(rc2)
				if tc.expectedB2 == nil {
					t.Errorf("expected rc2 to be nil; got %q", b2)
				} else if !bytes.Equal(b2, tc.expectedB2) {
					t.Errorf("expected rc2 %q; got %q", tc.expectedB1, b2)
				}
			} else if tc.expectedB2 != nil {
				t.Errorf("expected rc2 %q; got nil", tc.expectedB2)
			}

			if tc.expectedErr != err {
				t.Errorf("expected %s; got %s", tc.expectedErr, err)
			}
		})
	}
}

type errClose struct{}

func (errClose) Read([]byte) (n int, err error) {
	return 0, io.EOF
}

func (errClose) Close() error {
	return io.ErrUnexpectedEOF
}

type errRead struct{}

func (errRead) Read([]byte) (n int, err error) {
	return 0, io.ErrNoProgress
}

func (errRead) Close() error {
	return nil
}

type errReadClose struct{}

func (errReadClose) Read([]byte) (n int, err error) {
	return 0, io.ErrNoProgress
}

func (errReadClose) Close() error {
	return io.ErrUnexpectedEOF
}
