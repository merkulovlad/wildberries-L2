package scheduler

import (
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/merkulovlad/wildberries-L3/notification_service/internal/notification"
)

type WorkerScheduler struct {
	scheduler gocron.Scheduler
	worker    *notification.Worker
	interval  time.Duration
}

func NewWorkerScheduler(worker *notification.Worker, interval time.Duration) (*WorkerScheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &WorkerScheduler{
		scheduler: scheduler,
		worker:    worker,
		interval:  interval,
	}, nil
}

func (s *WorkerScheduler) Start(ctx context.Context) error {
	if _, err := s.scheduler.NewJob(
		gocron.DurationJob(s.interval),
		gocron.NewTask(func() error {
			return s.worker.RunOnce(ctx, time.Now())
		}),
	); err != nil {
		return err
	}

	s.scheduler.Start()
	return nil
}

func (s *WorkerScheduler) Shutdown() error {
	return s.scheduler.Shutdown()
}
