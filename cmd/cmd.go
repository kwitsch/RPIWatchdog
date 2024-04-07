package cmd

import (
	"context"

	"github.com/kwitsch/RPIWatchdog/logger"
)

const (
	OptionWatch       = "watch"
	OptionHealthcheck = "healthcheck"
	OptionSourcecheck = "sourcecheck"
	WrongOptionExit   = 1
)

func Run(option string) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch option {
	case OptionWatch:
		return Watch(ctx)
	case OptionHealthcheck:
		return Healthcheck(ctx)
	case OptionSourcecheck:
		return Sourcecheck(ctx)
	case "":
		logger.Log("No option provided")
		Help()
	default:
		logger.Log("Invalid option: %s", option)
		Help()
	}

	return WrongOptionExit
}

func Help() {
	logger.Log("Usage: rpiwatchdog [option]")
	logger.Log("Options:")
	logger.Log("\t- %s : Activates the watchdog device and maintains it's state acording to the source health\n", OptionWatch)
	logger.Log("\t- %s : Check the health of the service\n", OptionHealthcheck)
	logger.Log("\t- %s : Check the health of the source server\n", OptionSourcecheck)
}
