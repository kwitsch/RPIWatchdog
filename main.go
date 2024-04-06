package main

import (
	"errors"
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

	_, err = os.Stat(cfg.DevicePath)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatal("Configured watchdog device does not exist")
		os.Exit(2)
	}

	wd, err := watchdog.Open(cfg.DevicePath)
	if err != nil {
		log.Printf("Error opening watchdog: %v", err)
		os.Exit(3)
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
