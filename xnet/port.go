// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xnet

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
)

// FreePort asks the kernel for a free open port, that is ready to use, on the specified Network.
// Only TCP or UDP networks are supported.
func FreePort(ctx context.Context, network string, options ...ListenConfigOption) (int, error) {
	var lc net.ListenConfig

	for _, option := range options {
		option.apply(&lc)
	}

	switch network {
	case NetworkTCP, NetworkTCP4, NetworkTCP6:
		listener, err := lc.Listen(ctx, network, "localhost:0")
		if err != nil {
			return 0, err
		}
		defer listener.Close()
		return listener.Addr().(*net.TCPAddr).Port, nil
	case NetworkUDP, NetworkUDP4, NetworkUDP6:
		listener, err := lc.ListenPacket(ctx, network, "localhost:0")
		if err != nil {
			return 0, err
		}
		defer listener.Close()
		return listener.LocalAddr().(*net.UDPAddr).Port, nil
	default:
		return 0, fmt.Errorf("invalid network: %s", network)
	}
}

// ParsePort parses a string representing a port.
// If the string is not a valid port number, an error is returned.
func ParsePort(port string, allowZero bool) (int, error) {
	p, err := strconv.ParseUint(port, 10, 16) //nolint:gomnd // Base 10 number in range [0, 2^16-1]
	if err != nil {
		return 0, err
	}
	if p == 0 && !allowZero {
		return 0, errors.New("0 is not a valid port number")
	}
	return int(p), nil
}
