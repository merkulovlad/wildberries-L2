package notification

import (
	"context"
	"time"
)

type Worker struct {
	repo      Repository
	publisher Publisher
	limit     int
}

func NewWorker(repo Repository, publisher Publisher, limit int) *Worker {
	return &Worker{
		repo:      repo,
		publisher: publisher,
		limit:     limit,
	}
}

func (w *Worker) RunOnce(ctx context.Context, now time.Time) error {
	notifications, err := w.repo.ClaimDue(ctx, now, w.limit)
	if err != nil {
		return err
	}

	for _, n := range notifications {
		if err := w.publisher.Publish(ctx, n.ID); err != nil {
			return err
		}
	}

	return nil
}
