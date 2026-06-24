package notification

import (
	"context"
	"strings"
)

type Service interface {
	Create(ctx context.Context, in CreateNotification) (Notification, error)
	Get(ctx context.Context, id string) (Notification, error)
	Cancel(ctx context.Context, id string) (Notification, error)
}

type NotificationService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &NotificationService{repo: repo}
}

func (n *NotificationService) Create(ctx context.Context, in CreateNotification) (Notification, error) {
	in.Recipient = strings.TrimSpace(in.Recipient)
	if in.Recipient == "" {
		return Notification{}, ErrInvalidRecipient
	}

	in.Message = strings.TrimSpace(in.Message)
	if in.Message == "" {
		return Notification{}, ErrInvalidMessage
	}

	if in.SendAt.IsZero() {
		return Notification{}, ErrInvalidSendAt
	}

	return n.repo.Create(ctx, in)
}

func (n *NotificationService) Get(ctx context.Context, id string) (Notification, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Notification{}, ErrInvalidID
	}

	return n.repo.Get(ctx, id)
}

func (n *NotificationService) Cancel(ctx context.Context, id string) (Notification, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Notification{}, ErrInvalidID
	}

	if err := n.repo.Cancel(ctx, id); err != nil {
		return Notification{}, err
	}

	return n.repo.Get(ctx, id)
}
