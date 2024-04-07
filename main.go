package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
	_ "time/tzdata"

	"github.com/kwitsch/RPIWatchdog/config"
	"github.com/kwitsch/RPIWatchdog/healthcheck"

	_ "github.com/kwitsch/go-dockerutils"
)

const (
	optionWatch       = "watch"
	optionHealthcheck = "healthcheck"
	optionSourcecheck = "sourcecheck"
	wrongOptionExit   = 1
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Invalid number of arguments")
		help()
		os.Exit(wrongOptionExit)
	}
	switch os.Args[1] {
	case optionWatch:
		// todo
	case optionHealthcheck:
		// todo
	case optionSourcecheck:
	// todo
	default:
		log.Fatalf("Invalid option: %s", os.Args[1])
		help()
		os.Exit(wrongOptionExit)
	}
	// Get the configuration and exit with 1 if an error occurs
	cfg, err := config.GetConfig()
	if err != nil {
		log.Printf("Error getting config: %v", err)
		os.Exit(1)
	}

	// Open the watchdog device and exit with the corresponding error code(2|3) if an error occurs
	// wd, rcode := watchdog.Open(cfg.DevicePath)
	// if rcode != 0 {
	// 	os.Exit(rcode)
	// }
	// defer wd.Close()

	srv, err := healthcheck.NewHealthCheckServer(cfg.ServeHealthSource)
	if err != nil {
		log.Printf("Error creating health check server: %v", err)
		os.Exit(5)
	}
	srv.Serve(context.Background())
	defer srv.Close()

	// Listen for interrupts to terminate the service
	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, os.Interrupt)
	defer close(intChan)

	// Create ticker with 15 seconds interval for health checks
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-intChan:
			// wd.MagicClose()
			os.Exit(0)
		case <-srv.Err():
			log.Printf("Error serving health check: %v", <-srv.Err())
			os.Exit(7)
		case <-ticker.C:
			if err := healthcheck.UnixHealthCheck(context.Background()); err != nil {
				// if err := healthcheck.TCPHealthCheck(context.Background(), "127.0.0.1:1111"); err != nil {
				log.Printf("Error checking health: %v", err)
				os.Exit(6)
			} else {
				log.Printf("Health check successful")
			}
			// if err := wd.KeepAlive(); err != nil {
			// 	log.Printf("Error keeping watchdog alive: %v", err)
			// 	os.Exit(4)
			// }
		}
	}
}

func help() {
	log.Println("Usage: rpiwatchdog [option]")
	log.Println("Options:")
	log.Printf("\t- %s : Activates the watchdog device and maintains it's state acording to the source health\n", optionWatch)
	log.Printf("\t- %s : Check the health of the service\n", optionHealthcheck)
	log.Printf("\t- %s : Check the health of the source server\n", optionSourcecheck)
}
