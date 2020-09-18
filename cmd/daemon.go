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
	proxy  string
	access string
	secret string
)

func init() {
	rootCmd.AddCommand(daemonCmd)

	daemonCmd.Flags().StringVar(&proxy, "proxy", "", "address of relay-proxy for TCP connection")
	daemonCmd.MarkFlagRequired("proxy")

	daemonCmd.Flags().StringVar(&access, "access", "", "access key (aka username) to access the shell")
	daemonCmd.MarkFlagRequired("access")

	daemonCmd.Flags().StringVar(&secret, "secret", "", "secret key (aka password) to access the shell")
	daemonCmd.MarkFlagRequired("secret")
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "starts shell daemon and configures acceptance of connection via TCP proxy component",
	Long:  `starts shell daemon and configures acceptance of connection via TCP proxy component`,
	Example: `
sideshell daemon --proxy 127.0.0.1:8080 --access username --secret ABCD...EFGH
	`,
	SilenceUsage: true,
	RunE:         daemon,
}

func daemon(cmd *cobra.Command, args []string) error {
	sshd.ViaProxy(
		proxy,
		sshd.Credentials(access, secret),
		sshd.PrivateKey(4096),
	)
	return nil
}
