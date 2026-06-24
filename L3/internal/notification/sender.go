package notification

import "context"

type Sender interface {
	Send(ctx context.Context, n Notification) error
}
