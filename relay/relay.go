//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/sideshell
//

package proxy

import (
	"io"
	"log"
	"net"
)

//
func New(portA, portB string) {
	sideA := make(chan net.Conn)
	sideB := make(chan net.Conn)
	go listen(portA, sideA)
	go listen(portB, sideB)

	for tcpA := range sideA {
		select {
		case tcpB := <-sideB:
			go forward(tcpA, tcpB)
		}
	}
}

func listen(port string, router chan net.Conn) {
	log.Printf("Listening on %s...", port)
	listener, err := net.Listen("tcp", port)
	defer listener.Close()

	if err != nil {
		log.Fatalf("Failed to listen on %s (%s)", port, err)
	}

	for {
		tcp, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}
		log.Printf("Connected %s\n", tcp.RemoteAddr().String())
		router <- tcp
	}
}

func forward(tcpA, tcpB net.Conn) {
	copy := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			log.Printf("tcp copy error: %s\n", err)
			tcpA.Close()
			tcpB.Close()
			return
		}
	}

	log.Printf("forwarding %s <-> %s \n",
		tcpA.RemoteAddr().String(), tcpB.RemoteAddr().String())
	go copy(tcpA, tcpB)
	go copy(tcpB, tcpA)
}
