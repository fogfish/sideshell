//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/sideshell
//

package cmd

import (
	tcpRelay "github.com/fogfish/sideshell/relay"
	"github.com/spf13/cobra"
)

var (
	daemonPort string
	clientPort string
)

func init() {
	rootCmd.AddCommand(relayCmd)

	relayCmd.Flags().StringVar(&daemonPort, "daemon", "", "listen incoming daemon connection(s) on the port")
	relayCmd.MarkFlagRequired("daemon")

	relayCmd.Flags().StringVar(&clientPort, "client", "", "listen incoming client connection(s) on the port")
	relayCmd.MarkFlagRequired("client")
}

var relayCmd = &cobra.Command{
	Use:   "relay",
	Short: "relay TCP traffic between ports",
	Long:  `listens for incoming connection on ports and relays a traffic when connection is established`,
	Example: `
sideshell relay --daemon :8080 --client :8081
	`,
	SilenceUsage: true,
	RunE:         relay,
}

func relay(cmd *cobra.Command, args []string) error {
	tcpRelay.New(daemonPort, clientPort)
	return nil
}
