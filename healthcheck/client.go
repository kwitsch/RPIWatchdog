package healthcheck

import (
	"context"
	"fmt"
	"net"
	"time"
)

func UnixHealthCheck(ctx context.Context) error {
	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		return err
	}

	return healthcheck(ctx, conn)
}

func TCPHealthCheck(ctx context.Context, source string) error {
	if source == "" {
		return UnixHealthCheck(ctx)
	}
	conn, err := net.Dial("tcp", source)
	if err != nil {
		return err
	}

	return healthcheck(ctx, conn)
}

func healthcheck(ctx context.Context, conn net.Conn) error {
	defer conn.Close()
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
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
