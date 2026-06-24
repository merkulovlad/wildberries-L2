package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/merkulovlad/wildberries-L3/notification_service/internal/notification"
	wbredis "github.com/wb-go/wbf/redis"
)

var _ notification.Repository = (*Repository)(nil)

type Repository struct {
	next  notification.Repository
	redis *wbredis.Client
	ttl   time.Duration
}

func NewRepository(next notification.Repository, redis *wbredis.Client, ttl time.Duration) *Repository {
	return &Repository{
		next:  next,
		redis: redis,
		ttl:   ttl,
	}
}

func (r *Repository) Create(ctx context.Context, in notification.CreateNotification) (notification.Notification, error) {
	n, err := r.next.Create(ctx, in)
	if err != nil {
		return notification.Notification{}, err
	}
	_ = r.set(ctx, n)
	return n, nil
}

func (r *Repository) Get(ctx context.Context, id string) (notification.Notification, error) {
	n, err := r.get(ctx, id)
	if err == nil {
		return n, nil
	}

	n, err = r.next.Get(ctx, id)
	if err != nil {
		return notification.Notification{}, err
	}
	_ = r.set(ctx, n)
	return n, nil
}

func (r *Repository) Cancel(ctx context.Context, id string) error {
	if err := r.next.Cancel(ctx, id); err != nil {
		return err
	}
	_ = r.redis.Del(ctx, key(id))
	return nil
}

func (r *Repository) ClaimDue(ctx context.Context, now time.Time, limit int) ([]notification.Notification, error) {
	notifications, err := r.next.ClaimDue(ctx, now, limit)
	if err != nil {
		return nil, err
	}
	for _, n := range notifications {
		_ = r.set(ctx, n)
	}
	return notifications, nil
}

func (r *Repository) MarkSent(ctx context.Context, id string) error {
	if err := r.next.MarkSent(ctx, id); err != nil {
		return err
	}
	_ = r.redis.Del(ctx, key(id))
	return nil
}

func (r *Repository) MarkFailed(ctx context.Context, id string) error {
	if err := r.next.MarkFailed(ctx, id); err != nil {
		return err
	}
	_ = r.redis.Del(ctx, key(id))
	return nil
}

func (r *Repository) get(ctx context.Context, id string) (notification.Notification, error) {
	value, err := r.redis.Get(ctx, key(id))
	if err != nil {
		return notification.Notification{}, err
	}

	var n notification.Notification
	if err := json.Unmarshal([]byte(value), &n); err != nil {
		_ = r.redis.Del(ctx, key(id))
		return notification.Notification{}, err
	}

	return n, nil
}

func (r *Repository) set(ctx context.Context, n notification.Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return err
	}
	return r.redis.SetWithExpiration(ctx, key(n.ID), body, r.ttl)
}

func key(id string) string {
	return fmt.Sprintf("notification:%s", id)
}
