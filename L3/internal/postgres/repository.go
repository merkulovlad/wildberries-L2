package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/merkulovlad/wildberries-L3/notification_service/internal/notification"
	pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"
)

var _ notification.Repository = (*Repository)(nil)

type Repository struct {
	db pgxdriver.QueryExecuter
}

func NewRepository(db pgxdriver.QueryExecuter) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, in notification.CreateNotification) (notification.Notification, error) {
	const query = `
		INSERT INTO notifications (status, recipient, message, send_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, status, recipient, message, send_at
	`

	n, err := scanNotification(r.db.QueryRow(
		ctx,
		query,
		string(notification.StatusScheduled),
		in.Recipient,
		in.Message,
		in.SendAt,
	))
	if err != nil {
		return notification.Notification{}, err
	}

	return n, nil
}

func (r *Repository) Get(ctx context.Context, id string) (notification.Notification, error) {
	const query = `
		SELECT id, status, recipient, message, send_at
		FROM notifications
		WHERE id = $1
	`

	n, err := scanNotification(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return notification.Notification{}, notification.ErrNotFound
	}
	if err != nil {
		return notification.Notification{}, err
	}

	return n, nil
}

func (r *Repository) Cancel(ctx context.Context, id string) error {
	const query = `
		UPDATE notifications
		SET status = $2
		WHERE id = $1
			AND status IN ($3, $4)
	`

	tag, err := r.db.Exec(
		ctx,
		query,
		id,
		string(notification.StatusCanceled),
		string(notification.StatusScheduled),
		string(notification.StatusQueued),
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		current, getErr := r.Get(ctx, id)
		if getErr != nil {
			return getErr
		}
		if current.Status == notification.StatusCanceled {
			return notification.ErrAlreadyCanceled
		}
		return notification.ErrInvalidStatusTransition
	}

	return nil
}

func (r *Repository) ClaimDue(ctx context.Context, now time.Time, limit int) ([]notification.Notification, error) {
	if limit <= 0 {
		return nil, nil
	}

	const query = `
		WITH due AS (
			SELECT id
			FROM notifications
			WHERE status = $1
				AND send_at <= $2
			ORDER BY send_at
			LIMIT $3
			FOR UPDATE SKIP LOCKED
		)
		UPDATE notifications AS n
		SET status = $4
		FROM due
		WHERE n.id = due.id
		RETURNING n.id, n.status, n.recipient, n.message, n.send_at
	`

	rows, err := r.db.Query(
		ctx,
		query,
		string(notification.StatusScheduled),
		now,
		limit,
		string(notification.StatusQueued),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]notification.Notification, 0, limit)
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (r *Repository) MarkSent(ctx context.Context, id string) error {
	return r.updateStatus(ctx, id, notification.StatusSent)
}

func (r *Repository) MarkFailed(ctx context.Context, id string) error {
	return r.updateStatus(ctx, id, notification.StatusFailed)
}

func (r *Repository) updateStatus(ctx context.Context, id string, status notification.Status) error {
	const query = `
		UPDATE notifications
		SET status = $2
		WHERE id = $1
	`

	tag, err := r.db.Exec(ctx, query, id, string(status))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return notification.ErrNotFound
	}

	return nil
}

type notificationScanner interface {
	Scan(dest ...any) error
}

func scanNotification(row notificationScanner) (notification.Notification, error) {
	var n notification.Notification
	var status string

	err := row.Scan(&n.ID, &status, &n.Recipient, &n.Message, &n.SendAt)
	if err != nil {
		return notification.Notification{}, err
	}

	n.Status = notification.Status(status)
	return n, nil
}
