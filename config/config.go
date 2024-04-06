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
