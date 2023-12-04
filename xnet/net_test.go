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

func TestConn_Read(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(t *testing.T) (net.Listener, net.Conn)
		expectedErr bool
	}{
		{
			name: "connection already closed",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				t.Helper()
				ln, conn := dialTCPWithReadHandler(t, xnet.DialReadTimeout(5*time.Second))
				if err := conn.Close(); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return ln, conn
			},
			expectedErr: true,
		},
		{
			name: "no read timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				t.Helper()
				return dialTCPWithReadHandler(t)
			},
			expectedErr: false,
		},
		{
			name: "with negative read timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				t.Helper()
				return dialTCPWithReadHandler(t, xnet.DialReadTimeout(-5*time.Second))
			},
			expectedErr: true,
		},
		{
			name: "with positive read timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				t.Helper()
				return dialTCPWithReadHandler(t, xnet.DialReadTimeout(5*time.Second))
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ln, conn := tc.setup(t)

			buf := make([]byte, 4)
			n, err := conn.Read(buf)

			conn.Close()
			ln.Close()

			assertOperation(t, tc.expectedErr, n, err)
		})
	}
}

func TestConn_Write(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(t *testing.T) (net.Listener, net.Conn)
		expectedErr bool
	}{
		{
			name: "connection already closed",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				t.Helper()
				ln, conn := dialTCPWithWriteHandler(t, xnet.DialWriteTimeout(5*time.Second))
				if err := conn.Close(); err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return ln, conn
			},
			expectedErr: true,
		},
		{
			name: "no write timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				t.Helper()
				return dialTCPWithWriteHandler(t)
			},
			expectedErr: false,
		},
		{
			name: "with negative write timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				t.Helper()
				return dialTCPWithWriteHandler(t, xnet.DialWriteTimeout(-5*time.Second))
			},
			expectedErr: true,
		},
		{
			name: "with positive write timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				t.Helper()
				return dialTCPWithWriteHandler(t, xnet.DialWriteTimeout(5*time.Second))
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ln, conn := tc.setup(t)

			n, err := conn.Write([]byte(`ping`))

			conn.Close()
			ln.Close()

			assertOperation(t, tc.expectedErr, n, err)
		})
	}
}

func assertOperation(tb testing.TB, expectedErr bool, n int, err error) {
	tb.Helper()

	isErrNil := err == nil
	if expectedErr == isErrNil {
		tb.Errorf("expected error is %t, got %v", expectedErr, err)
	}

	noBytes := n == 0
	if expectedErr != noBytes {
		tb.Errorf("expected bytes is %t, got %d bytes", !expectedErr, n)
	}
}

func handleConnections(tb testing.TB, ln net.Listener, handler func(net.Conn) error) {
	tb.Helper()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}

		if err := handler(conn); err != nil {
			tb.Logf("connection handling error: %s", err)
		}
		if err := conn.Close(); err != nil {
			tb.Logf("connection closure error: %s", err)
		}
	}
}

func dialTCP(tb testing.TB, handler func(net.Conn) error, options ...xnet.DialOption) (net.Listener, net.Conn, error) {
	tb.Helper()

	ln, port, err := listenTCP()
	if err != nil {
		return nil, nil, err
	}

	go handleConnections(tb, ln, handler)

	conn, err := xnet.DialContext(context.Background(), xnet.NetworkTCP, net.JoinHostPort("127.0.0.1", port), options...)
	if err != nil {
		if cerr := ln.Close(); cerr != nil {
			return nil, nil, errors.Join(err, cerr)
		}
		return nil, nil, err
	}

	return ln, conn, nil
}

func dialTCPWithReadHandler(tb testing.TB, options ...xnet.DialOption) (net.Listener, net.Conn) {
	tb.Helper()

	ln, conn, err := dialTCP(tb, func(conn net.Conn) error {
		_, err := conn.Write([]byte("pong"))
		return err
	}, options...)
	if err != nil {
		tb.Fatal(err)
	}
	return ln, conn
}

func dialTCPWithWriteHandler(tb testing.TB, options ...xnet.DialOption) (net.Listener, net.Conn) {
	tb.Helper()

	ln, conn, err := dialTCP(tb, func(conn net.Conn) error {
		_, err := conn.Read(make([]byte, 4))
		return err
	}, options...)
	if err != nil {
		tb.Fatal(err)
	}
	return ln, conn
}
