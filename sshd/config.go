//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/sideshell
//

package sshd

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

//
type Config func(*Server)

//
func Credentials(access, secret string) Config {
	return func(s *Server) {
		s.config.PasswordCallback = func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in a production setting.
			if c.User() == access && string(pass) == secret {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		}
	}
}

//
func PrivateKey(bits int) Config {
	return func(s *Server) {
		key, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			log.Fatalf("Failed to generate private key (%d bits)", bits)
		}

		err = key.Validate()
		if err != nil {
			log.Fatalf("Failed to validate private key: %v", err)
		}

		signer, err := ssh.NewSignerFromKey(key)
		if err != nil {
			log.Fatalf("Failed to create private key: %v", err)
		}

		s.config.AddHostKey(signer)
	}
}

//
func PrivateKeyFile(idfile string) Config {
	return func(s *Server) {
		bytes, err := ioutil.ReadFile("id_rsa")
		if err != nil {
			log.Fatalf("Failed to load private key (%s)", idfile)
		}

		key, err := ssh.ParsePrivateKey(bytes)
		if err != nil {
			log.Fatal("Failed to parse private key")
		}
		s.config.AddHostKey(key)
	}
}

//
func Shell(shell string) Config {
	return func(s *Server) {
		s.shell = shell
	}
}
