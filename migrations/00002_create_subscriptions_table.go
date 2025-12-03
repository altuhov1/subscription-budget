package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateSubscriptionsTable, downCreateSubscriptionsTable)
}

func upCreateSubscriptionsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE subscriptions (
			id SERIAL PRIMARY KEY,
			service_name VARCHAR(255) NOT NULL,
			price INTEGER NOT NULL CHECK (price > 0),
			user_id UUID NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE
		);
	`)
	if err != nil {
		return nil
	}
	_, err = tx.ExecContext(ctx, `
		CREATE INDEX idx_subscriptions_query_optimized 
		ON subscriptions (user_id, start_date, end_date, service_name);
	`)
	return err
}

func downCreateSubscriptionsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		DROP TABLE IF EXISTS subscriptions CASCADE;
	`)
	return err
}
