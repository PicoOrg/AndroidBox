package main

import (
	"github.com/PicoOrg/AndroidBox/internal/util"
	"github.com/spf13/cobra"
)

var (
	logger  util.Logger
	verbose bool
)

var (
	rootCmd = &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger = util.NewLogger(verbose)
		},
	}
)

func init() {
	rootCmd.AddCommand(&cobra.Command{Use: "completion", Hidden: true})
	rootCmd.AddCommand(setpropCmd)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose")
}
