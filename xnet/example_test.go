// Copyright 2023 Jérémy Lourenço. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xnet_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jlourenc/xgo/xnet"
)

func ExampleDial() {
	dialOptions := []xnet.DialOption{
		xnet.DialReadTimeout(time.Second),
		xnet.DialWriteTimeout(time.Second),
	}

	conn, err := xnet.Dial(xnet.NetworkTCP, "localhost:12345", dialOptions...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	log.Print("Connection established")
}

func ExampleDialContext() {
	dialOptions := []xnet.DialOption{
		xnet.DialReadTimeout(time.Second),
		xnet.DialWriteTimeout(time.Second),
	}

	conn, err := xnet.DialContext(context.Background(), xnet.NetworkTCP, "localhost:12345", dialOptions...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	log.Print("Connection established")
}

func ExampleDialer_Dial() {
	d := xnet.Dialer{
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	conn, err := d.Dial(xnet.NetworkTCP, "localhost:12345")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	log.Print("Connection established")
}

func ExampleDialer_DialContext() {
	d := xnet.Dialer{
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	conn, err := d.DialContext(context.Background(), xnet.NetworkTCP, "localhost:12345")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	log.Print("Connection established")
}

func ExampleConn_Read() {
	conn, err := xnet.Dial(xnet.NetworkTCP, "localhost:12345")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	if _, err := conn.Read(buffer); err != nil {
		log.Print(err)
	}
}

func ExampleConn_Write() {
	conn, err := xnet.Dial(xnet.NetworkTCP, "localhost:12345")
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("Hello, World!")); err != nil {
		log.Print(err)
	}
}

func ExampleFreePort() {
	port, err := xnet.FreePort(context.Background(), xnet.NetworkTCP)
	if err != nil {
		log.Fatalf("Failed to get free port: %v", err)
	}

	log.Printf("port: %d", port)
}

func ExampleParsePort() {
	port, err := xnet.ParsePort("12001", false)
	if err != nil {
		log.Fatalf("Failed to parse port number: %v", err)
	}

	fmt.Printf("%d\n", port)
	// Output: 12001
}
