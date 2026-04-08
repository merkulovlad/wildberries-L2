package shell

import (
	"os"
	"os/signal"
)

type SignalHandler struct {
	signals chan os.Signal
}

// NewSignalHandler creates a new signal handler for managing OS signals.
func NewSignalHandler() *SignalHandler {
	return &SignalHandler{
		signals: make(chan os.Signal, 1),
	}
}

// SetupWithInterrupt initializes signal handling for SIGINT.
func (h *SignalHandler) SetupWithInterrupt() {
	signal.Notify(h.signals, os.Interrupt)
}

// C returns the signal channel for long-lived listeners.
func (h *SignalHandler) C() <-chan os.Signal {
	return h.signals
}

func (h *SignalHandler) Stop() {
	signal.Stop(h.signals)
}
