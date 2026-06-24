package notification

import "context"

type Publisher interface {
	Publish(ctx context.Context, notificationID string) error
}
