package main

import (
	"os"
	_ "time/tzdata"

	"github.com/kwitsch/RPIWatchdog/cmd"

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

	os.Exit(cmd.Run(option))
}
