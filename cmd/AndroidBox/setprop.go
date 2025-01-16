package main

import (
	"github.com/PicoOrg/AndroidBox/internal/mprop"
	"github.com/PicoOrg/AndroidBox/internal/util"
	"github.com/spf13/cobra"
)

var (
	initPid int
)

var (
	setpropCmd = &cobra.Command{
		Use:   "setprop name value",
		Short: "setprop",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				cmd.Help()
				return
			}
			mprop := mprop.NewMProp(logger, initPid)
			name, value := args[0], args[1]
			err := mprop.Set(name, value)
			if err != nil {
				logger.Info("mprop set failed, please retry", util.Fields{"error": err, "name": name, "value": value})
			} else {
				logger.Info("mprop set success", util.Fields{"name": name, "value": value})
			}
		},
	}
)

func init() {
	setpropCmd.PersistentFlags().IntVarP(&initPid, "initPid", "p", 1, "init pid")
}
