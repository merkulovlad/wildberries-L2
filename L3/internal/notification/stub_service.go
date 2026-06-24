package notification

import "context"

type notImplementedService struct{}

func NewNotImplementedService() Service {
	return notImplementedService{}
}

func (notImplementedService) Create(context.Context, CreateNotification) (Notification, error) {
	return Notification{}, ErrNotImplemented
}

func (notImplementedService) Get(context.Context, string) (Notification, error) {
	return Notification{}, ErrNotImplemented
}

func (notImplementedService) Cancel(context.Context, string) (Notification, error) {
	return Notification{}, ErrNotImplemented
}
