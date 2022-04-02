// Copyright 2022 Jérémy Lourenço. All rights reserved.
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

func listenTCP() (net.Listener, string, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, "", err
	}
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		ln.Close()
		return nil, "", err
	}
	return ln, port, nil
}

func assertDial(t *testing.T, expectedErr bool, conn net.Conn, err error) {
	t.Helper()

	if expectedErr {
		if conn != nil {
			t.Errorf("expected no connection, got %v", conn)
		}
		if err == nil {
			t.Error("expected error, got nil")
		}
	} else {
		if conn == nil {
			t.Error("expected connection, got nil")
		}
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
	}
}

func TestDial(t *testing.T) {
	testCases := []struct {
		name        string
		network     string
		options     []DialOption
		expectedErr bool
	}{
		{
			name:        "invalid network",
			network:     "invalid",
			expectedErr: true,
		},
		{
			name:    "valid network",
			network: NetworkTCP,
			options: []DialOption{
				DialConnectDeadline(time.Now().Add(5 * time.Second)),
				DialConnectTimeout(5 * time.Second),
				DialKeepAlive(30 * time.Second),
				DialReadTimeout(5 * time.Second),
				DialWriteTimeout(5 * time.Second),
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

			conn, err := Dial(tc.network, net.JoinHostPort("127.0.0.1", port), tc.options...)

			assertDial(t, tc.expectedErr, conn, err)
		})
	}
}

func TestDialContext(t *testing.T) {
	testCases := []struct {
		name        string
		network     string
		options     []DialOption
		expectedErr bool
	}{
		{
			name:        "invalid network",
			network:     "invalid",
			expectedErr: true,
		},
		{
			name:    "valid network",
			network: NetworkTCP,
			options: []DialOption{
				DialConnectDeadline(time.Now().Add(5 * time.Second)),
				DialConnectTimeout(5 * time.Second),
				DialKeepAlive(30 * time.Second),
				DialReadTimeout(5 * time.Second),
				DialWriteTimeout(5 * time.Second),
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

			conn, err := DialContext(context.Background(), tc.network, net.JoinHostPort("127.0.0.1", port), tc.options...)

			assertDial(t, tc.expectedErr, conn, err)
		})
	}
}
