package notification

import (
	"errors"
	"time"
)

var (
	ErrNotImplemented          = errors.New("notification service is not implemented")
	ErrNotFound                = errors.New("notification not found")
	ErrInvalidID               = errors.New("notification id is required")
	ErrInvalidRecipient        = errors.New("notification recipient is required")
	ErrInvalidMessage          = errors.New("notification message is required")
	ErrInvalidSendAt           = errors.New("notification send_at is required")
	ErrAlreadyCanceled         = errors.New("notification already canceled")
	ErrInvalidStatusTransition = errors.New("invalid notification status transition")
)

type Notification struct {
	ID        string
	Status    Status
	Recipient string
	Message   string
	SendAt    time.Time
}

type CreateNotification struct {
	Recipient string
	Message   string
	SendAt    time.Time
}

type Status string

const (
	StatusScheduled Status = "scheduled"
	StatusQueued    Status = "queued"
	StatusSent      Status = "sent"
	StatusCanceled  Status = "canceled"
	StatusFailed    Status = "failed"
)
