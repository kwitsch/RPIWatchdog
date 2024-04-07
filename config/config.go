package config

import (
	. "github.com/kwitsch/go-dockerutils/config"
)

// prefix is the prefix for the environment variables
const prefix = "RPIW_"

// Config is the configuration for the watchdog container
type Config struct {
	// DevicePath is the path to the watchdog device to use
	DevicePath string `koanf:"devicepath" default:"/dev/watchdog"`
	// ServeHealthSource determines if the health endpoint should be exposed on port 1111
	ServeHealthSource bool `koanf:"servehealthsource" default:"false"`
	// UseHealthSource is the hostname and port to use as the health source
	// If it is empty, the health source will be the internal unix socket
	UseHealthSource string `koanf:"usehealthsource" default:""`
}

// GetConfig returns the configuration for the watchdog container according to the environment variables and
// docker secrets
func GetConfig() (Config, error) {
	var res Config
	if err := Load(prefix, &res); err != nil {
		return res, err
	}

	return res, nil
}
