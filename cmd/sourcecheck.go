package cmd

import (
	"context"

	"github.com/kwitsch/RPIWatchdog/config"
	"github.com/kwitsch/RPIWatchdog/healthcheck"
	"github.com/kwitsch/RPIWatchdog/logger"
)

// Sourcecheck checks the health of the source and displays the result
func Sourcecheck(ctx context.Context) int {
	cfg, res := config.GetConfig()
	if res != 0 {
		logger.Log("Error reading config")
		return res
	}

	if err := healthcheck.TCPHealthCheck(ctx, cfg.UseHealthSource, cfg.HealthCheckTimeout); err != nil {
		logger.Log("Health check for %s was failed: %v", cfg.UseHealthSource, err)
		return 1
	}

	logger.Log("Health check for %s was successful", cfg.UseHealthSource)

	return 0
}
