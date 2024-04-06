package main

import (
	"os"
	"os/signal"
	_ "time/tzdata"

	_ "github.com/kwitsch/go-dockerutils"
)

func main() {
	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, os.Interrupt)
	defer close(intChan)

	for {
		select {
		case <-intChan:
			os.Exit(0)
		}
	}
}
