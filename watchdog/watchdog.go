package watchdog

import (
	"errors"
	"os"
	"time"

	"github.com/kwitsch/RPIWatchdog/logger"
	wd "github.com/mdlayher/watchdog"
)

const (
	DevicePath             = "/dev/watchdog"
	defaultTimeout         = 30 * time.Second
	deviceDoesNotExistExit = 11
	watchdogOpenErrorExit  = 12
)

type Watchdog struct {
	wdDevice        *wd.Device
	withoutWatchdog bool
	timeout         time.Duration
}

// NewWatchdog creates a new watchdog instance and returns it
// If the watchdog device does not exist or can't be opened, it returns an exit code
func NewWatchdog(withoutWatchdog bool) (Watchdog, int) {
	res := Watchdog{
		withoutWatchdog: withoutWatchdog,
	}
	if withoutWatchdog {
		logger.Log("!!!!!!!!!!!!!!!!!!!!!!!!!!")
		logger.Log("!! Watchdog is disabled !!")
		logger.Log("!!!!!!!!!!!!!!!!!!!!!!!!!!")

		return res, 0
	}

	if _, err := os.Stat(DevicePath); errors.Is(err, os.ErrNotExist) {
		logger.Log("Configured watchdog device does not exist")
		return res, deviceDoesNotExistExit
	}

	wd, err := wd.Open()
	if err != nil {
		logger.Log("Error opening watchdog: %v", err)
		return res, watchdogOpenErrorExit
	}
	res.wdDevice = wd

	timeout, err := wd.Timeout()
	if err != nil {
		logger.Log("Error getting watchdog timeout: %v", err)
		wd.Close()
		return res, watchdogOpenErrorExit
	}
	res.timeout = timeout

	logger.LogVerbose("Watchdog:")
	logger.LogVerbose(" - Driver: %s", wd.Identity)
	logger.LogVerbose(" - Timeout: %v", timeout)
	logger.LogVerbose("----------------------------------------------")

	return res, 0
}

// Timeout returns the timeout of the watchdog device
func (w *Watchdog) Timeout() time.Duration {
	if w.withoutWatchdog {
		return defaultTimeout
	}

	return w.timeout
}

// Ping sends a keep alive signal to the watchdog device
func (w *Watchdog) Ping() error {
	if w.withoutWatchdog {
		logger.Log("! Watchdog is disabled ! - keep alive not called")
		return nil
	}

	logger.LogVerbose("Watchdog - keep alive called")
	return w.wdDevice.Ping()
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
