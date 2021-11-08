// Copyright 2021 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xnet extends the Go standard library package net
// by providing additional primitives and structures for network I/O.
package xnet

import (
	"net"
	"time"
)

// Enumaration of networks.
const (
	NetworkIP         = "ip"
	NetworkIP4        = "ip4"
	NetworkIP6        = "ip6"
	NetworkTCP        = "tcp"
	NetworkTCP4       = "tcp4"
	NetworkTCP6       = "tcp6"
	NetworkUDP        = "udp"
	NetworkUDP4       = "udp4"
	NetworkUDP6       = "udp6"
	NetworkUnix       = "unix"
	NetworkUnixgram   = "unixgram"
	NetworkUnixpacket = "unixpacket"
)

type conn struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// Read reads data from the connection.
// No deadline is set if the Conn read timeout is the zero value.
// A deadline, defined as current time + read timeout, is set otherwise.
//
// See net.Conn.Read for more information.
func (c *conn) Read(b []byte) (int, error) {
	if c.readTimeout != 0 {
		if err := c.Conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return 0, err
		}
	}
	return c.Conn.Read(b)
}

// Write writes data to the connection.
// No deadline is set if the Conn write timeout is the zero value.
// A deadline, defined as current time + write timeout, is set otherwise.
//
// See net.Conn.Write for more information.
func (c *conn) Write(b []byte) (int, error) {
	if c.writeTimeout != 0 {
		if err := c.Conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return 0, err
		}
	}
	return c.Conn.Write(b)
}
