package watchdog

import (
	"errors"
	"log"
	"os"

	wd "github.com/u-root/u-root/pkg/watchdog"
)

func Open(devicePath string) (*wd.Watchdog, int) {
	_, err := os.Stat(devicePath)
	if errors.Is(err, os.ErrNotExist) {
		log.Fatal("Configured watchdog device does not exist")
		return nil, 2
	}

	wd, err := wd.Open(devicePath)
	if err != nil {
		log.Printf("Error opening watchdog: %v", err)
		return nil, 3
	}

	return wd, 0
}
