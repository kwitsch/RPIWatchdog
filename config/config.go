package config

import (
	"github.com/kwitsch/RPIWatchdog/logger"
	. "github.com/kwitsch/go-dockerutils/config"
)

const (
	// prefix is the prefix for the environment variables
	prefix = "RPIW_"
	// configErrorExit is the exit code if an error occurs while getting the configuration
	configErrorExit = 2
)

// Config is the configuration for the watchdog container
type Config struct {
	// DevicePath is the path to the watchdog device to use
	DevicePath string `koanf:"devicepath" default:"/dev/watchdog"`
	// ServeHealthSource determines if the health endpoint should be exposed on port 1111
	ServeHealthSource bool `koanf:"servehealthsource" default:"false"`
	// UseHealthSource is the hostname and port to use as the health source
	// If it is empty, the health source will be the internal unix socket
	UseHealthSource string `koanf:"usehealthsource" default:""`
	// VerboseLogging is the flag to enable or disable additional log messages
	VerboseLogging bool `koanf:"verboselogging" default:"false"`
	// HealthCheckInterval is the interval in seconds to check the health source
	HealthCheckInterval int `koanf:"healthcheckinterval" default:"30"`
	// HealthCheckTimeout is the timeout in seconds to check the health source
	HealthCheckTimeout int `koanf:"healthchecktimeout" default:"3"`
}

// GetConfig returns the configuration for the watchdog container according to the environment variables and
// docker secrets
func GetConfig() (Config, int) {
	var res Config
	if err := Load(prefix, &res); err != nil {
		logger.Log("Error loading config: %v", err)
		return res, configErrorExit
	}

	logger.SetVerbose(res.VerboseLogging)

	return res, 0
}
