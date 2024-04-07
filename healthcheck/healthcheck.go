package healthcheck

const (
	// tcpPort is the port to expose the health endpoint on if ExposeHealth is set to true
	tcpPort = 1111
	// sockPath is the path to the health socket
	sockPath = "/tmp/watchdog-health.sock"
	// healthMessage is the message to send for a successful health check
	healthMessage = "OK"
)
