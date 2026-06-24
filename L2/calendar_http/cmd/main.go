package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/merkulovlad/wildberries-L2/calendar_http/internal/calendar"
	"github.com/merkulovlad/wildberries-L2/calendar_http/internal/httpapi"
)

func main() {
	addr := flag.String("addr", envOr("CALENDAR_ADDR", ":8080"), "listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "calendar ", log.LstdFlags|log.Lmicroseconds)

	cal := calendar.New()
	srv := httpapi.NewServer(cal)
	handler := httpapi.Logging(logger, srv.Handler())

	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		logger.Fatalf("listen: %v", err)
	}
	logger.Printf("listening on %s", ln.Addr())

	errCh := make(chan error, 1)
	go func() { errCh <- httpSrv.Serve(ln) }()

	select {
	case <-ctx.Done():
		logger.Printf("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(shutdownCtx); err != nil {
			logger.Printf("shutdown: %v", err)
		}
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("serve: %v", err)
		}
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
