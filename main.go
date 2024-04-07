package main

import (
	"os"
	_ "time/tzdata"

	"github.com/kwitsch/RPIWatchdog/cmd"
	"github.com/kwitsch/RPIWatchdog/logger"

	reaper "github.com/ramr/go-reaper"
)

func init() {
	go reaper.Reap()
}

func main() {
	option := ""
	if len(os.Args) > 1 {
		option = os.Args[1]
	}

	exitCode := cmd.Run(option)

	logger.LogVerbose("Exit code: %d", exitCode)

	os.Exit(exitCode)
}
