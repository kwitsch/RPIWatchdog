package watchdog

import (
	"errors"
	"os"

	"github.com/kwitsch/RPIWatchdog/logger"
	wd "github.com/u-root/u-root/pkg/watchdog"
)

const (
	deviceDoesNotExistExit = 10
	watchdogOpenErrorExit  = 11
)

type Watchdog struct {
	wdDevice        *wd.Watchdog
	withoutWatchdog bool
}

// NewWatchdog creates a new watchdog instance and returns it
// If the watchdog device does not exist or can't be opened, it returns an exit code
func NewWatchdog(devicePath string, withoutWatchdog bool) (Watchdog, int) {
	res := Watchdog{
		withoutWatchdog: withoutWatchdog,
	}
	if withoutWatchdog {
		logger.Log("!!!!!!!!!!!!!!!!!!!!!!!!!!")
		logger.Log("!! Watchdog is disabled !!")
		logger.Log("!!!!!!!!!!!!!!!!!!!!!!!!!!")

		return res, 0
	}

	if _, err := os.Stat(devicePath); errors.Is(err, os.ErrNotExist) {
		logger.Log("Configured watchdog device does not exist")
		return res, deviceDoesNotExistExit
	}

	wd, err := wd.Open(devicePath)
	if err != nil {
		logger.Log("Error opening watchdog: %v", err)
		return res, watchdogOpenErrorExit
	}
	res.wdDevice = wd

	return res, 0
}

// KeepAlive sends a keep alive signal to the watchdog device
func (w *Watchdog) KeepAlive() error {
	if w.withoutWatchdog {
		logger.Log("! Watchdog is disabled ! - keep alive not called")
		return nil
	}

	logger.LogVerbose("Watchdog - keep alive called")
	return w.wdDevice.KeepAlive()
}

// Close closes the watchdog device
func (w *Watchdog) Close() error {
	if w.withoutWatchdog {
		logger.Log("! Watchdog is disabled ! - close not called")
		return nil
	}

	logger.LogVerbose("Watchdog - close called")
	return w.wdDevice.Close()
}

// Disable disables the watchdog device
func (w *Watchdog) Disable() error {
	if w.withoutWatchdog {
		logger.Log("! Watchdog is disabled ! - disable not called")
		return nil
	}

	logger.LogVerbose("Watchdog - disable called")
	return w.wdDevice.MagicClose()
}
