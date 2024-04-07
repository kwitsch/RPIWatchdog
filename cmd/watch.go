package cmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/kwitsch/RPIWatchdog/config"
	"github.com/kwitsch/RPIWatchdog/healthcheck"
	"github.com/kwitsch/RPIWatchdog/logger"
	"github.com/kwitsch/RPIWatchdog/watchdog"
)

const (
	healthcheckServerCreateErrorExit = 5
	healthCheckServerRunErrorExit    = 6
	contextDoneExit                  = 7
	watchdogCreateErrorExit          = 8
	watchdogKeepAliveErrorExit       = 9
)

// Watch starts the RPIWatchdog service
func Watch(ctx context.Context) int {
	logger.Log("Start RPIWatchdog service")

	cfg, res := config.GetConfig()
	if res != 0 {
		logger.Log("Error reading config")
		return res
	}

	srv, err := healthcheck.NewHealthCheckServer(cfg.ServeHealthSource, cfg.HealthCheckTimeout)
	if err != nil {
		logger.Log("Error creating health check server: %v", err)
		return healthcheckServerCreateErrorExit
	}
	srv.Serve(context.Background())
	defer srv.Close()

	// Listen for interrupts to terminate the service
	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, os.Interrupt)
	defer close(intChan)

	// Create ticker with 15 seconds interval for health checks
	ticker := time.NewTicker(time.Duration(cfg.HealthCheckInterval) * time.Second)
	defer ticker.Stop()

	// Open the watchdog device and exit with the corresponding error code(2|3) if an error occurs
	wd, rcode := watchdog.NewWatchdog(cfg.DevicePath, cfg.WithoutWatchdog)
	if rcode != 0 {
		return watchdogCreateErrorExit
	}
	defer wd.Close()

	if res := checkHealth(ctx, cfg, &wd); res != 0 {
		return res
	}

	for {
		select {
		case <-ctx.Done():
			logger.LogVerbose("Context done")
			return contextDoneExit
		case <-intChan:
			logger.LogVerbose("Interrupt signal received")
			wd.Disable()
			return 0
		case err = <-srv.Err():
			logger.Log("Error serving health check: %v", err)
			return healthCheckServerRunErrorExit
		case <-ticker.C:
			if res := checkHealth(ctx, cfg, &wd); res != 0 {
				return res
			}
		}
	}
}

func checkHealth(ctx context.Context, cfg config.Config, wd *watchdog.Watchdog) int {
	logger.LogVerbose("Health check")

	logSource := cfg.DevicePath
	if cfg.UseHealthSource != "" {
		logSource = cfg.UseHealthSource
	}

	if err := healthcheck.TCPHealthCheck(ctx, cfg.UseHealthSource, cfg.HealthCheckTimeout); err != nil {
		logger.LogVerbose(" - failed for %s with: %v", logSource, err)
	} else {
		logger.LogVerbose(" - was successful for %s", logSource)
		if err := wd.KeepAlive(); err != nil {
			logger.Log("Error keeping watchdog alive: %v", err)
			return watchdogKeepAliveErrorExit
		}
	}

	return 0
}
