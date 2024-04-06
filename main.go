package main

import (
	"log"
	"os"
	"os/signal"
	_ "time/tzdata"

	"github.com/kwitsch/RPIWatchdog/config"
	"github.com/u-root/u-root/pkg/watchdog"

	_ "github.com/kwitsch/go-dockerutils"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Printf("Error getting config: %v", err)
		os.Exit(1)
	}

	wd, err := watchdog.Open(cfg.DevicePath)
	if err != nil {
		log.Printf("Error opening watchdog: %v", err)
		os.Exit(2)
	}
	defer wd.Close()

	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, os.Interrupt)
	defer close(intChan)

	for {
		select {
		case <-intChan:
			wd.MagicClose()
			os.Exit(0)
		default:
			if err := wd.KeepAlive(); err != nil {
				log.Printf("Error keeping watchdog alive: %v", err)
				os.Exit(4)
			}
		}
	}
}
