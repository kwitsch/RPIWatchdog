package logger

import (
	"fmt"
	"log"
)

// verbose is the flag to enable or disable verbose logging
var verbose bool = false

// SetVerbose sets the verbose flag to enable or disable verbose logging
func SetVerbose(v bool) {
	verbose = v
}

// Log writes a log message as line to the standard logger after formating it
func Log(format string, v ...any) {
	log.Println(fmt.Sprintf(format, v...))
}

// LogVerbose writes a log message as line to the standard logger after formating it if the verbose flag is set
func LogVerbose(format string, v ...any) {
	if verbose {
		log.Println(fmt.Sprintf(format, v...))
	}
}
