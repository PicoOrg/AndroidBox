package main

import (
	"os"

	"github.com/PicoOrg/AndroidBox/internal/mprop"
	"github.com/PicoOrg/AndroidBox/internal/util"
)

const DEBUG = false

func main() {
	logger := util.NewLogger(DEBUG)
	mprop := mprop.NewMProp(logger, 1)
	err := mprop.Set(os.Args[1], os.Args[2])
	if err != nil {
		logger.Info("mprop set failed, please retry", util.Fields{"error": err, "name": os.Args[1], "value": os.Args[2]})
	} else {
		logger.Info("mprop set success", util.Fields{"name": os.Args[1], "value": os.Args[2]})
	}
}
