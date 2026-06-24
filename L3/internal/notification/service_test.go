package notification

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestServiceCreate(t *testing.T) {
	sendAt := time.Date(2026, 6, 24, 15, 0, 0, 0, time.UTC)
	repo := fakeRepository{
		createFn: func(_ context.Context, in CreateNotification) (Notification, error) {
			if in.Recipient != "user-1" {
				t.Fatalf("recipient = %q, want user-1", in.Recipient)
			}
			if in.Message != "hello" {
				t.Fatalf("message = %q, want hello", in.Message)
			}
			if !in.SendAt.Equal(sendAt) {
				t.Fatalf("send_at = %s, want %s", in.SendAt, sendAt)
			}
			return Notification{ID: "n-1", Status: StatusScheduled}, nil
		},
	}

	got, err := NewService(repo).Create(context.Background(), CreateNotification{
		Recipient: " user-1 ",
		Message:   " hello ",
		SendAt:    sendAt,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if got.ID != "n-1" || got.Status != StatusScheduled {
		t.Fatalf("Create() = %+v, want created notification", got)
	}
}

func TestServiceCreateValidation(t *testing.T) {
	tests := []struct {
		name string
		in   CreateNotification
		want error
	}{
		{name: "recipient", in: CreateNotification{Message: "hello", SendAt: time.Now()}, want: ErrInvalidRecipient},
		{name: "message", in: CreateNotification{Recipient: "user-1", SendAt: time.Now()}, want: ErrInvalidMessage},
		{name: "send_at", in: CreateNotification{Recipient: "user-1", Message: "hello"}, want: ErrInvalidSendAt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewService(fakeRepository{}).Create(context.Background(), tt.in)
			if !errors.Is(err, tt.want) {
				t.Fatalf("Create() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestServiceGetValidation(t *testing.T) {
	_, err := NewService(fakeRepository{}).Get(context.Background(), " ")
	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("Get() error = %v, want %v", err, ErrInvalidID)
	}
}

func TestServiceCancel(t *testing.T) {
	repo := fakeRepository{
		cancelFn: func(_ context.Context, id string) error {
			if id != "n-1" {
				t.Fatalf("cancel id = %q, want n-1", id)
			}
			return nil
		},
		getFn: func(_ context.Context, id string) (Notification, error) {
			if id != "n-1" {
				t.Fatalf("get id = %q, want n-1", id)
			}
			return Notification{ID: id, Status: StatusCanceled}, nil
		},
	}

	got, err := NewService(repo).Cancel(context.Background(), " n-1 ")
	if err != nil {
		t.Fatalf("Cancel() error = %v", err)
	}
	if got.ID != "n-1" || got.Status != StatusCanceled {
		t.Fatalf("Cancel() = %+v, want canceled notification", got)
	}
}

type fakeRepository struct {
	createFn   func(context.Context, CreateNotification) (Notification, error)
	getFn      func(context.Context, string) (Notification, error)
	cancelFn   func(context.Context, string) error
	claimDueFn func(context.Context, time.Time, int) ([]Notification, error)
}

func (r fakeRepository) Create(ctx context.Context, in CreateNotification) (Notification, error) {
	if r.createFn == nil {
		return Notification{}, ErrNotImplemented
	}
	return r.createFn(ctx, in)
}

func (r fakeRepository) Get(ctx context.Context, id string) (Notification, error) {
	if r.getFn == nil {
		return Notification{}, ErrNotImplemented
	}
	return r.getFn(ctx, id)
}

func (r fakeRepository) Cancel(ctx context.Context, id string) error {
	if r.cancelFn == nil {
		return ErrNotImplemented
	}
	return r.cancelFn(ctx, id)
}

func (r fakeRepository) ClaimDue(ctx context.Context, now time.Time, limit int) ([]Notification, error) {
	if r.claimDueFn == nil {
		return nil, ErrNotImplemented
	}
	return r.claimDueFn(ctx, now, limit)
}

func (r fakeRepository) MarkSent(context.Context, string) error {
	return ErrNotImplemented
}

func (r fakeRepository) MarkFailed(context.Context, string) error {
	return ErrNotImplemented
}
