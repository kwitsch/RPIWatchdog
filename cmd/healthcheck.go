package cmd

import (
	"context"

	"github.com/kwitsch/RPIWatchdog/config"
	"github.com/kwitsch/RPIWatchdog/healthcheck"
	"github.com/kwitsch/RPIWatchdog/logger"
)

const (
	successExit = 0
	errorExit   = 1
)

func Healthcheck(ctx context.Context) int {
	cfg, res := config.GetConfig()
	if res != 0 {
		logger.Log("Error reading config")
		return errorExit
	}

	if err := healthcheck.UnixHealthCheck(ctx, cfg.HealthCheckTimeout); err != nil {
		logger.LogVerbose("Health check failed: %v", err)
		return errorExit
	}

	logger.LogVerbose("Health check successful")

	return successExit
}
