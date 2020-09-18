//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/sideshell
//

package main

import (
	"fmt"
	"os"

	"github.com/fogfish/sideshell/sshd"
)

func main() {
	fmt.Println("==> Hello World!!!")

	sshd.ViaProxy(
		os.Getenv("CONFIG_SSHD_PROXY"),
		sshd.Credentials(
			os.Getenv("CONFIG_SSHD_ACCESS"),
			os.Getenv("CONFIG_SSHD_SECRET"),
		),
		sshd.PrivateKey(4096),
	)
}
