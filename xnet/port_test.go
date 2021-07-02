// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xnet_test

import (
	"context"
	"errors"
	"syscall"
	"testing"
	"time"

	. "github.com/jlourenc/xgo/xnet"
)

func TestFreePort(t *testing.T) {
	testCases := []struct {
		name        string
		options     []ListenConfigOption
		network     string
		expectedErr bool
	}{
		{
			name:        "unsupported network",
			network:     NetworkUnix,
			expectedErr: true,
		},
		{
			name:        "tcp network - success",
			network:     NetworkTCP,
			expectedErr: false,
		},
		{
			name:        "udp network - success",
			network:     NetworkUDP,
			expectedErr: false,
		},
		{
			name: "tcp network - failure",
			options: []ListenConfigOption{
				ListenConfigControl(func(network, address string, c syscall.RawConn) error {
					return errors.New("always error")
				}),
				ListenConfigKeepAlive(time.Second),
			},
			network:     NetworkTCP,
			expectedErr: true,
		},
		{
			name: "udp network - failure",
			options: []ListenConfigOption{
				ListenConfigControl(func(network, address string, c syscall.RawConn) error {
					return errors.New("always error")
				}),
				ListenConfigKeepAlive(time.Second),
			},
			network:     NetworkUDP,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			port, err := FreePort(context.Background(), tc.network, tc.options...)

			if tc.expectedErr {
				if port != 0 {
					t.Errorf("expected 0, got %d", port)
				}
				if err == nil {
					t.Error("expected an error, got nil")
				}
			} else {
				if port == 0 {
					t.Error("expected a port, got 0")
				}
				if err != nil {
					t.Errorf("expected no error, got %s", err)
				}
			}
		})
	}
}
