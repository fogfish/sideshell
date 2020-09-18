//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/sideshell
//

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "sideshell",
	Short:   "xxx",
	Long:    `xxx`,
	Run:     root,
	Version: "v0",
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}
