//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/sideshell
//

package sshd

import (
	"io"
	"log"
	"os"
	"os/exec"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func osShellSession(shell string, channel ssh.Channel, control <-chan *ssh.Request) {
	stream := terminal.NewTerminal(channel, "")

	cmd := exec.Command(shell, "-i")
	cmd.Env = os.Environ()

	shing, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("failed to set input %v", err)
	}

	shout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("failed to set output %v", err)
	}

	sherr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("failed to set error %v", err)
	}

	go func() { io.Copy(stream, shout) }()
	go func() { io.Copy(stream, sherr) }()
	go func() { io.Copy(shing, channel) }()

	err = cmd.Run()
	if err != nil {
		log.Printf("failed to run %v", err)
	}
	log.Printf("Closing session ")
	channel.Close()
}
