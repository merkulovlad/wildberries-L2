package notification

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, in CreateNotification) (Notification, error)
	Get(ctx context.Context, id string) (Notification, error)
	Cancel(ctx context.Context, id string) error
	ClaimDue(ctx context.Context, now time.Time, limit int) ([]Notification, error)
	MarkSent(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string) error
}
