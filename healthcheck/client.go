package healthcheck

import (
	"context"
	"fmt"
	"net"
	"time"
)

// UnixHealthCheck checks the health through the unix socket
func UnixHealthCheck(ctx context.Context, timeout int) error {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		return err
	}

	return healthcheck(ctx, conn, timeout)
}

// TCPHealthCheck checks the health through the tcp health source if it is not empty, otherwise through the unix socket
func TCPHealthCheck(ctx context.Context, source string, timeout int) error {
	if source == "" {
		return UnixHealthCheck(ctx, timeout)
	}
	conn, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}

	return healthcheck(ctx, conn, timeout)
}

// healthcheck checks the health of the connection
func healthcheck(ctx context.Context, conn net.Conn, timeout int) error {
	defer conn.Close()
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeout))); err != nil {
			return err
		}
		buf := make([]byte, 2)
		_, err := conn.Read(buf)
		if err != nil {
			return err
		}

		if res := string(buf); res != healthMessage {
			return fmt.Errorf("Health check failed(Received: %s, Expected: OK)", res)
		}
	}
	return nil
}
