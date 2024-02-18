// Copyright 2024 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xnet_test

import (
	"context"
	"errors"
	"syscall"
	"testing"
	"time"

	"github.com/jlourenc/xgo/xnet"
)

func TestFreePort(t *testing.T) {
	testCases := []struct {
		name        string
		options     []xnet.ListenConfigOption
		network     string
		expectedErr bool
	}{
		{
			name:        "unsupported network",
			network:     xnet.NetworkUnix,
			expectedErr: true,
		},
		{
			name:        "tcp network - success",
			network:     xnet.NetworkTCP,
			expectedErr: false,
		},
		{
			name:        "udp network - success",
			network:     xnet.NetworkUDP,
			expectedErr: false,
		},
		{
			name: "tcp network - failure",
			options: []xnet.ListenConfigOption{
				xnet.ListenConfigControl(func(network, address string, c syscall.RawConn) error {
					return errors.New("always error")
				}),
				xnet.ListenConfigKeepAlive(time.Second),
			},
			network:     xnet.NetworkTCP,
			expectedErr: true,
		},
		{
			name: "udp network - failure",
			options: []xnet.ListenConfigOption{
				xnet.ListenConfigControl(func(network, address string, c syscall.RawConn) error {
					return errors.New("always error")
				}),
				xnet.ListenConfigKeepAlive(time.Second),
			},
			network:     xnet.NetworkUDP,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			port, err := xnet.FreePort(context.Background(), tc.network, tc.options...)

			isErrNil := err == nil
			if tc.expectedErr == isErrNil {
				t.Errorf("expected error is %t, got %v", tc.expectedErr, err)
			}

			noPortNumber := port == 0
			if tc.expectedErr != noPortNumber {
				t.Errorf("expected port is %t, got %d", !tc.expectedErr, port)
			}
		})
	}
}

func TestParsePort(t *testing.T) {
	testCases := []struct {
		name         string
		port         string
		allowZero    bool
		expectedPort int
		expectedErr  bool
	}{
		{
			name:         "port not a number",
			port:         "invalid",
			allowZero:    false,
			expectedPort: 0,
			expectedErr:  true,
		},
		{
			name:         "negative port number",
			port:         "-1",
			allowZero:    false,
			expectedPort: 0,
			expectedErr:  true,
		},
		{
			name:         "port number greater than 2^16-1",
			port:         "65536",
			allowZero:    false,
			expectedPort: 0,
			expectedErr:  true,
		},
		{
			name:         "port number is 0 and 0 is not allowed",
			port:         "0",
			allowZero:    false,
			expectedPort: 0,
			expectedErr:  true,
		},
		{
			name:         "port number is 0 and 0 is allowed",
			port:         "0",
			allowZero:    true,
			expectedPort: 0,
			expectedErr:  false,
		},
		{
			name:         "port number is valid",
			port:         "12001",
			allowZero:    false,
			expectedPort: 12001,
			expectedErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			port, err := xnet.ParsePort(tc.port, tc.allowZero)

			if tc.expectedPort != port {
				t.Errorf("expected %d, got %d", tc.expectedPort, port)
			}
			if tc.expectedErr {
				if err == nil {
					t.Error("expected an error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %s", err)
				}
			}
		})
	}
}
