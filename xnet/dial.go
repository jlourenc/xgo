// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xnet

import (
	"context"
	"net"
	"syscall"
	"time"
)

// Dialer is a wrapper around net.Dialer that provides additional options for connecting to an address.
//
// See net.Dialer for more information.
type Dialer struct {
	net.Dialer
	// ReadTimeout is the maximum amount of time after which each Conn Read operation
	// fails instead of blocking. In this case, an error that wraps os.ErrDeadlineExceeded
	// is returned. This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
	//
	// The error's Timeout method will return true, but note that there are other possible
	// errors for which the Timeout method will return true even if the Read has timed out.
	//
	// The default is no timeout. (zero value)
	ReadTimeout time.Duration
	// WriteTimeout is the maximum amount of time after which each Conn Write operation
	// fails instead of blocking. In this case, an error that wraps os.ErrDeadlineExceeded
	// is returned. This can be tested using errors.Is(err, os.ErrDeadlineExceeded).
	//
	// The error's Timeout method will return true, but note that there are other possible
	// errors for which the Timeout method will return true even if the Write has timed out.
	//
	// The default is no timeout. (zero value)
	WriteTimeout time.Duration
}

// Dial acts like net.Dial but uses a Dialer that supports read and write timeouts at the connection level.
// Optional DialOption parameters may be passed in to configure the Dialer.
//
// Read and write timeouts are respectively applied to each Read and Write call on the net.Conn (a zero value means no timeout).
//
// See net.Dial for more information.
func Dial(network, address string, options ...DialOption) (net.Conn, error) {
	var d Dialer
	for _, option := range options {
		option.apply(&d)
	}
	return d.Dial(network, address)
}

// DialContext acts like Dial but takes a context.Context.
//
// See net.Dialer.DialContext for more information.
func DialContext(ctx context.Context, network, address string, options ...DialOption) (net.Conn, error) {
	var d Dialer
	for _, option := range options {
		option.apply(&d)
	}
	return d.DialContext(ctx, network, address)
}

// Dial acts like net.Dialer.Dial but uses a Dialer that supports read and write timeouts at the connection level.
//
// See Dial for more information.
func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

// DialContext acts like Dialer.Dial but takes a context.Context.
//
// See net.Dialer.DialContext for more information.
func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	c, err := d.Dialer.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	return &conn{Conn: c, readTimeout: d.ReadTimeout, writeTimeout: d.WriteTimeout}, nil
}

type (
	// DialOption configures how connections are made.
	DialOption interface {
		apply(d *Dialer)
	}

	funcDialOption struct {
		fn func(*Dialer)
	}
)

func newFuncDialOption(fn func(*Dialer)) funcDialOption {
	return funcDialOption{
		fn: fn,
	}
}

func (fdo funcDialOption) apply(d *Dialer) {
	fdo.fn(d)
}

// DialConnectDeadline returns a DialOption that configures a deadline for a connect to complete by.
func DialConnectDeadline(deadline time.Time) DialOption {
	return newFuncDialOption(func(d *Dialer) {
		d.Deadline = deadline
	})
}

// DialConnectTimeout returns a DialOption that configures a timeout for a connect to complete.
func DialConnectTimeout(timeout time.Duration) DialOption {
	return newFuncDialOption(func(d *Dialer) {
		d.Timeout = timeout
	})
}

// WithKeepAlive returns a DialOption that configures the interval
// between keep-alive probes for an active network TCP connection.
func DialKeepAlive(keepAlive time.Duration) DialOption {
	return newFuncDialOption(func(d *Dialer) {
		d.KeepAlive = keepAlive
	})
}

// DialReadTimeout returns a DialOption that configures a timeout for a Conn Read to complete.
func DialReadTimeout(timeout time.Duration) DialOption {
	return newFuncDialOption(func(d *Dialer) {
		d.ReadTimeout = timeout
	})
}

// DialWriteTimeout returns a DialOption that configures a timeout for a Conn Write to complete.
func DialWriteTimeout(timeout time.Duration) DialOption {
	return newFuncDialOption(func(d *Dialer) {
		d.WriteTimeout = timeout
	})
}

type (
	// ListenConfigOption is an option to configure  anet.ListenConfig.
	ListenConfigOption interface {
		apply(lc *net.ListenConfig)
	}

	funcListenConfigOption struct {
		fn func(*net.ListenConfig)
	}
)

func newFuncListenConfigOption(fn func(*net.ListenConfig)) funcListenConfigOption {
	return funcListenConfigOption{
		fn: fn,
	}
}

func (flco funcListenConfigOption) apply(lc *net.ListenConfig) {
	flco.fn(lc)
}

// ListenConfigControl returns a ListenConfigOption that configures a control of a net.ListenConfig.
func ListenConfigControl(fn func(network, address string, c syscall.RawConn) error) ListenConfigOption {
	return newFuncListenConfigOption(func(lc *net.ListenConfig) {
		lc.Control = fn
	})
}

// ListenConfigKeepAlive returns a ListenConfigOption that configures a keep alive for a net.ListenConfig.
func ListenConfigKeepAlive(keepAlive time.Duration) ListenConfigOption {
	return newFuncListenConfigOption(func(lc *net.ListenConfig) {
		lc.KeepAlive = keepAlive
	})
}
