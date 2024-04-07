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
	// ServeHealthSource determines if the health endpoint should be exposed on port 1111
	ServeHealthSource bool `koanf:"servehealthsource" default:"false"`
	// UseHealthSource is the hostname and port to use as the health source
	// If it is empty, the health source will be the internal unix socket
	UseHealthSource string `koanf:"usehealthsource" default:""`
	// VerboseLogging is the flag to enable or disable additional log messages
	VerboseLogging bool `koanf:"verboselogging" default:"false"`
	// HealthCheckTimeout is the timeout in seconds to check the health source
	HealthCheckTimeout int `koanf:"healthchecktimeout" default:"3"`
	// WithoutWatchdog is the flag to disable the watchdog device
	// If it is true, the watchdog device will not be opened
	// This flag is only for testing purposes
	WithoutWatchdog bool `koanf:"withoutwatchdog" default:"false"`
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

	logConfig(res)

	return res, 0
}

func logConfig(cfg Config) {
	logger.LogVerbose("----------------------------------------------")
	logger.LogVerbose("Configuration:")
	logger.LogVerbose(" - ServeHealthSource: %t", cfg.ServeHealthSource)
	logger.LogVerbose(" - UseHealthSource: %s", cfg.UseHealthSource)
	logger.LogVerbose(" - VerboseLogging: %t", cfg.VerboseLogging)
	logger.LogVerbose(" - HealthCheckTimeout: %d", cfg.HealthCheckTimeout)
	logger.LogVerbose(" - WithoutWatchdog: %t", cfg.WithoutWatchdog)
	logger.LogVerbose("----------------------------------------------")
}
