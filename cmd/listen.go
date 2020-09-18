//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/sideshell
//

package cmd

import (
	"github.com/fogfish/sideshell/sshd"
	"github.com/spf13/cobra"
)

var (
	port string
)

func init() {
	rootCmd.AddCommand(listenCmd)

	listenCmd.Flags().StringVar(&port, "port", "", "listen incoming connection on the port")
	listenCmd.MarkFlagRequired("port")

	listenCmd.Flags().StringVar(&access, "access", "", "access key (aka username) to access the shell")
	listenCmd.MarkFlagRequired("access")

	listenCmd.Flags().StringVar(&secret, "secret", "", "secret key (aka password) to access the shell")
	listenCmd.MarkFlagRequired("secret")
}

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "listen for shell sessions on the port",
	Long:  `listen for shell sessions on the port`,
	Example: `
sideshell listen --port :8080 --access username --secret ABCD...EFGH
	`,
	SilenceUsage: true,
	RunE:         listen,
}

func listen(cmd *cobra.Command, args []string) error {
	sshd.Listen(
		port,
		sshd.Credentials(access, secret),
		sshd.PrivateKey(4096),
	)
	return nil
}
