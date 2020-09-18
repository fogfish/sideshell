//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/sideshell
//

// Inspired by
// https://gist.github.com/jpillora/b480fde82bff51a06238

package sshd

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

//
type server struct {
	config *ssh.ServerConfig
	shell  string
}

// ViaProxy listen ssh session using tcp relay-proxy
func ViaProxy(host string, opts ...Config) {
	daemon := New(opts...)

	for {
		log.Printf("Connecting %s...", host)

		tcp, err := net.Dial("tcp", host)
		if err != nil {
			log.Printf("Failed to connect (%s)", err)
			time.Sleep(5 * time.Second)
			continue
		}

		daemon.Accept(tcp)
	}

}

// Listen ssh sessions on the port
func Listen(port string, opts ...Config) {
	daemon := New(opts...)

	log.Printf("Listening on %s...", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %s (%s)", port, err)
	}

	for {
		tcp, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}

		daemon.Accept(tcp)
	}
}

//
func New(opts ...Config) *server {
	srv := &server{
		config: &ssh.ServerConfig{},
		shell:  "sh",
	}
	for _, opt := range opts {
		opt(srv)
	}
	return srv
}

//
func (srv *server) Accept(sock net.Conn) error {
	// New creates ssh daemon on the given connection
	sshSock, chans, reqs, err := ssh.NewServerConn(sock, srv.config)

	if err != nil {
		log.Printf("Failed to handshake (%s)", err)
		return err
	}

	log.Printf("New SSH connection from %s (%s)",
		sshSock.RemoteAddr(), sshSock.ClientVersion())

	// Discard all global out-of-band Requests
	go ssh.DiscardRequests(reqs)
	// Accept all channels
	go srv.handleChannels(chans)

	return nil
}

func (srv *server) handleChannels(chans <-chan ssh.NewChannel) {
	// Service the incoming Channel channel in go routine
	for newChannel := range chans {
		go srv.handleChannel(newChannel)
	}
}

func (srv *server) handleChannel(newChannel ssh.NewChannel) {
	// Since we're handling a shell, we expect a
	// channel type of "session". The also describes
	// "x11", "direct-tcpip" and "forwarded-tcpip"
	// channel types.
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}

	// At this point, we have the opportunity to reject the client's
	// request for another logical connection
	channel, control, err := newChannel.Accept()
	if err != nil {
		log.Printf("Could not accept channel (%s)", err)
		return
	}

	osShellSession(srv.shell, channel, control)
}
