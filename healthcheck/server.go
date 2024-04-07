package healthcheck

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
)

// HealthCheckServer is the server for the health checks on unix and tcp sockets
type HealthCheckServer struct {
	unixSocket   net.Listener
	tcpSocket    net.Listener
	errCh        chan error
	writeTimeout int
}

// NewHealthCheckServer creates a new health check server
func NewHealthCheckServer(serveHealthSource bool, writeTimeout int) (HealthCheckServer, error) {
	res := HealthCheckServer{
		errCh:        make(chan error, 1),
		writeTimeout: writeTimeout,
	}

	// Remove the unix socket if it already exists
	if err := os.Remove(sockPath); err != nil && !os.IsNotExist(err) {
		return res, err
	}

	// Create the unix socket
	unixSocket, err := net.Listen("unix", sockPath)
	if err != nil {
		return res, err
	}
	res.unixSocket = unixSocket

	// Create the tcp socket if exposehealth is set to true
	if serveHealthSource {
		tcpSocket, err := net.Listen("tcp", fmt.Sprintf(":%d", tcpPort))
		if err != nil {
			res.unixSocket.Close()
			return res, err
		}
		res.tcpSocket = tcpSocket
	}

	return res, nil
}

// Serve starts the health check server sockets
func (h *HealthCheckServer) Serve(ctx context.Context) {
	go func(gctx context.Context) {
		h.errCh <- h.serveUnix(gctx)
	}(ctx)
	if h.tcpSocket != nil {
		go func(gctx context.Context) {
			h.errCh <- h.serveTCP(gctx)
		}(ctx)
	}
}

// Err returns the error channel
func (h *HealthCheckServer) Err() <-chan error {
	return h.errCh
}

// Close closes the health check server
func (h *HealthCheckServer) Close() {
	if h.unixSocket != nil {
		h.unixSocket.Close()
	}
	if h.tcpSocket != nil {
		h.tcpSocket.Close()
	}
}

// serveUnix serves the health check on the unix socket by writing "OK" to the connection
func (h *HealthCheckServer) serveUnix(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := h.unixSocket.Accept()
			if err != nil {
				return err
			}
			go h.handleConnection(ctx, conn)
		}
	}
}

// serveTCP serves the health check on the tcp socket by writing "OK" to the connection
func (h *HealthCheckServer) serveTCP(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			conn, err := h.tcpSocket.Accept()
			if err != nil {
				return err
			}
			go h.handleConnection(ctx, conn)
		}
	}
}

// handleConnection writes "OK" to the connection and closes it
func (h *HealthCheckServer) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	select {
	case <-ctx.Done():
		return
	default:
		if err := conn.SetWriteDeadline(time.Now().Add(time.Duration(h.writeTimeout) * time.Second)); err != nil {
			h.errCh <- err
			return
		}
		if _, err := conn.Write([]byte(healthMessage)); err != nil {
			h.errCh <- err
			return
		}
	}
}
