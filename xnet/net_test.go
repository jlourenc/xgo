// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xnet_test

import (
	"context"
	"net"
	"testing"
	"time"

	. "github.com/jlourenc/xgo/xnet"
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
				ln, conn := dialTCPWithReadHandler(t, DialReadTimeout(5*time.Second))
				conn.Close()
				return ln, conn
			},
			expectedErr: true,
		},
		{
			name: "no read timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				return dialTCPWithReadHandler(t)
			},
			expectedErr: false,
		},
		{
			name: "with negative read timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				return dialTCPWithReadHandler(t, DialReadTimeout(-5*time.Second))
			},
			expectedErr: true,
		},
		{
			name: "with positive read timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				return dialTCPWithReadHandler(t, DialReadTimeout(5*time.Second))
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
				ln, conn := dialTCPWithWriteHandler(t, DialWriteTimeout(5*time.Second))
				conn.Close()
				return ln, conn
			},
			expectedErr: true,
		},
		{
			name: "no write timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				return dialTCPWithWriteHandler(t)
			},
			expectedErr: false,
		},
		{
			name: "with negative write timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				return dialTCPWithWriteHandler(t, DialWriteTimeout(-5*time.Second))
			},
			expectedErr: true,
		},
		{
			name: "with positive write timeout",
			setup: func(t *testing.T) (net.Listener, net.Conn) {
				return dialTCPWithWriteHandler(t, DialWriteTimeout(5*time.Second))
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

func assertOperation(t *testing.T, expectedErr bool, n int, err error) {
	if expectedErr {
		if n != 0 {
			t.Errorf("expected no bytes, got %d bytes", n)
		}
		if err == nil {
			t.Error("expected error, got nil")
		}
	} else {
		if n == 0 {
			t.Error("expected bytes, got no bytes")
		}
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
	}
}

func handleConnections(ln net.Listener, handler func(net.Conn) error) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}

		handler(conn)
		conn.Close()
	}
}

func dialTCP(handler func(net.Conn) error, options ...DialOption) (net.Listener, net.Conn, error) {
	ln, port, err := listenTCP()
	if err != nil {
		return nil, nil, err
	}

	go handleConnections(ln, handler)

	conn, err := DialContext(context.Background(), NetworkTCP, net.JoinHostPort("127.0.0.1", port), options...)
	if err != nil {
		ln.Close()
		return nil, nil, err
	}

	return ln, conn, nil
}

func dialTCPWithReadHandler(t *testing.T, options ...DialOption) (net.Listener, net.Conn) {
	ln, conn, err := dialTCP(func(conn net.Conn) error {
		_, err := conn.Write([]byte("pong"))
		return err
	}, options...)
	if err != nil {
		t.Fatal(err)
	}
	return ln, conn
}

func dialTCPWithWriteHandler(t *testing.T, options ...DialOption) (net.Listener, net.Conn) {
	ln, conn, err := dialTCP(func(conn net.Conn) error {
		_, err := conn.Read(make([]byte, 4))
		return err
	}, options...)
	if err != nil {
		t.Fatal(err)
	}
	return ln, conn
}
