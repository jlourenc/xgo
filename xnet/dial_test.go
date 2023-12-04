// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xnet_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/jlourenc/xgo/xnet"
)

func listenTCP() (net.Listener, string, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, "", err
	}
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		if cerr := ln.Close(); cerr != nil {
			return nil, "", errors.Join(err, cerr)
		}
		return nil, "", err
	}
	return ln, port, nil
}

func assertDial(tb testing.TB, expectedErr bool, conn net.Conn, err error) {
	tb.Helper()

	isErrNil := err == nil
	if expectedErr == isErrNil {
		tb.Errorf("expected error is %t, got %v", expectedErr, err)
	}

	isConnNil := conn == nil
	if expectedErr != isConnNil {
		tb.Errorf("expected connection is %t, got %v", !expectedErr, conn)
	}
}

func TestDial(t *testing.T) {
	testCases := []struct {
		name        string
		network     string
		options     []xnet.DialOption
		expectedErr bool
	}{
		{
			name:        "invalid network",
			network:     "invalid",
			expectedErr: true,
		},
		{
			name:    "valid network",
			network: xnet.NetworkTCP,
			options: []xnet.DialOption{
				xnet.DialConnectDeadline(time.Now().Add(5 * time.Second)),
				xnet.DialConnectTimeout(5 * time.Second),
				xnet.DialKeepAlive(30 * time.Second),
				xnet.DialReadTimeout(5 * time.Second),
				xnet.DialWriteTimeout(5 * time.Second),
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ln, port, err := listenTCP()
			if err != nil {
				t.Fatal(err)
			}
			defer ln.Close()

			conn, err := xnet.Dial(tc.network, net.JoinHostPort("127.0.0.1", port), tc.options...)

			assertDial(t, tc.expectedErr, conn, err)
		})
	}
}

func TestDialContext(t *testing.T) {
	testCases := []struct {
		name        string
		network     string
		options     []xnet.DialOption
		expectedErr bool
	}{
		{
			name:        "invalid network",
			network:     "invalid",
			expectedErr: true,
		},
		{
			name:    "valid network",
			network: xnet.NetworkTCP,
			options: []xnet.DialOption{
				xnet.DialConnectDeadline(time.Now().Add(5 * time.Second)),
				xnet.DialConnectTimeout(5 * time.Second),
				xnet.DialKeepAlive(30 * time.Second),
				xnet.DialReadTimeout(5 * time.Second),
				xnet.DialWriteTimeout(5 * time.Second),
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ln, port, err := listenTCP()
			if err != nil {
				t.Fatal(err)
			}
			defer ln.Close()

			conn, err := xnet.DialContext(context.Background(), tc.network, net.JoinHostPort("127.0.0.1", port), tc.options...)

			assertDial(t, tc.expectedErr, conn, err)
		})
	}
}
