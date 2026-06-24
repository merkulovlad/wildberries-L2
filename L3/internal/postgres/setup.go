package postgres

import (
	"context"

	pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"
)

func Setup(ctx context.Context, db pgxdriver.QueryExecuter) error {
	const query = `
		CREATE EXTENSION IF NOT EXISTS pgcrypto;

		CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
			status TEXT NOT NULL,
			recipient TEXT NOT NULL,
			message TEXT NOT NULL,
			send_at TIMESTAMPTZ NOT NULL
		);

		CREATE INDEX IF NOT EXISTS notifications_due_idx
			ON notifications (status, send_at);
	`

	_, err := db.Exec(ctx, query)
	return err
}
