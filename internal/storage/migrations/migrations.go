package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUserAndTable(ctx context.Context, pool *pgxpool.Pool, username, password, tableName string) error {
	createUserSQL := fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s'", username, password)
	_, err := pool.Exec(ctx, createUserSQL)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	createTableSQL := `
        CREATE TABLE subscriptions (
            id SERIAL PRIMARY KEY,
            service_name VARCHAR(255) NOT NULL,
            price INTEGER NOT NULL CHECK (price > 0),
            user_id UUID NOT NULL,
            start_date DATE NOT NULL,
            end_date DATE
        )`

	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	grantSQL := fmt.Sprintf("GRANT ALL PRIVILEGES ON TABLE subscriptions TO %s", username)
	_, err = pool.Exec(ctx, grantSQL)
	if err != nil {
		return fmt.Errorf("failed to grant privileges: %w", err)
	}

	grantSeqSQL := fmt.Sprintf("GRANT ALL PRIVILEGES ON SEQUENCE subscriptions_id_seq TO %s", username)
	_, err = pool.Exec(ctx, grantSeqSQL)
	if err != nil {
		return fmt.Errorf("failed to grant sequence privileges: %w", err)
	}

	return nil
}
