package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

// Run dials host:port with the given timeout, then proxies data between
// the connection and stdin/stdout until either side closes or ctx is done.
func Run(ctx context.Context, host, port string, timeout time.Duration) error {
	addr := net.JoinHostPort(host, port)

	dialer := net.Dialer{Timeout: timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", addr, err)
	}
	defer conn.Close()

	fmt.Fprintf(os.Stderr, "connected to %s\n", addr)

	return pipe(ctx, conn)
}

// pipe runs two goroutines proxying conn <-> stdin/stdout. It returns when
// the server closes the connection, the user presses Ctrl+D, or ctx is done.
func pipe(ctx context.Context, conn net.Conn) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		mu    sync.Mutex
		rxErr error
	)

	// conn -> stdout: terminates when the server closes the connection.
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil && !errors.Is(err, net.ErrClosed) {
			mu.Lock()
			rxErr = err
			mu.Unlock()
		}
		cancel()
	}()

	// stdin -> conn: terminates on Ctrl+D (EOF on stdin).
	go func() {
		_, _ = io.Copy(conn, os.Stdin)
		// Half-close the write side so the server sees EOF. The rx goroutine
		// finishes once the server closes its side. For non-TCP conns just
		// cancel.
		if tcp, ok := conn.(*net.TCPConn); ok {
			_ = tcp.CloseWrite()
		} else {
			cancel()
		}
	}()

	<-ctx.Done()
	_ = conn.Close()

	mu.Lock()
	defer mu.Unlock()
	return rxErr
}
