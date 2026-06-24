package notification

import (
	"context"
	"testing"
	"time"
)

func TestWorkerRunOnce(t *testing.T) {
	now := time.Date(2026, 6, 24, 15, 0, 0, 0, time.UTC)
	repo := fakeRepository{
		claimDueFn: func(_ context.Context, gotNow time.Time, limit int) ([]Notification, error) {
			if !gotNow.Equal(now) {
				t.Fatalf("now = %s, want %s", gotNow, now)
			}
			if limit != 10 {
				t.Fatalf("limit = %d, want 10", limit)
			}
			return []Notification{
				{ID: "n-1", Status: StatusQueued},
				{ID: "n-2", Status: StatusQueued},
			}, nil
		},
	}
	publisher := fakePublisher{}

	err := NewWorker(repo, &publisher, 10).RunOnce(context.Background(), now)
	if err != nil {
		t.Fatalf("RunOnce() error = %v", err)
	}

	if len(publisher.ids) != 2 || publisher.ids[0] != "n-1" || publisher.ids[1] != "n-2" {
		t.Fatalf("published ids = %v, want [n-1 n-2]", publisher.ids)
	}
}

type fakePublisher struct {
	ids []string
}

func (p *fakePublisher) Publish(_ context.Context, notificationID string) error {
	p.ids = append(p.ids, notificationID)
	return nil
}
