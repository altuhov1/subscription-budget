package storage

import (
	"context"
	"fmt"
	"time"

	"budget/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionStoragePG struct {
	pool *pgxpool.Pool
}

func NewSubscriptionStoragePG(pool *pgxpool.Pool) *SubscriptionStoragePG {
	return &SubscriptionStoragePG{pool: pool}
}

func (s *SubscriptionStoragePG) CreateSubscription(ctx context.Context, req *models.SubscriptionForStorage) (int, error) {
	var endDate pgtype.Date
	if req.EndDate != nil {
		endDate = pgtype.Date{Time: *req.EndDate, Valid: true}
	} else {
		endDate = pgtype.Date{Valid: false}
	}

	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int
	err := s.pool.QueryRow(ctx, query,
		req.ServiceName,
		req.Price,
		req.UserID,
		req.StartDate,
		endDate,
	).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("ошибка при вставке в БД: %w", err)
	}

	return id, nil
}

func (s *SubscriptionStoragePG) GetSubscriptionsWithParam(ctx context.Context, req *models.TotalCostRequestForStorage) (*models.TotalCostResponse, error) {
	var endDate pgtype.Date
	if req.EndDate != nil {
		endDate = pgtype.Date{Time: *req.EndDate, Valid: true}
	} else {
		endDate = pgtype.Date{Valid: false}
	}
	query := `
		SELECT COALESCE(SUM(
			price * (
				EXTRACT(YEAR FROM AGE(LEAST(end_date, $3), GREATEST(start_date, $2))) * 12
				+ 1 + EXTRACT(MONTH FROM AGE(LEAST(end_date, $3), GREATEST(start_date, $2)))
			)
		), 0) AS total_cost
		FROM subscriptions
		WHERE ($1::text IS NULL OR $1 = '' OR user_id = $1::uuid)
		AND start_date <= $3
		AND (end_date IS NULL OR end_date >= $2)
		AND ($4::text[] IS NULL OR service_name = ANY($4))
	`

	var total int
	err := s.pool.QueryRow(ctx, query,
		req.UserID,
		req.StartDate,
		endDate,
		req.Subscriptions,
	).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("ошибка при расчёте общей стоимости: %w", err)
	}

	return &models.TotalCostResponse{Total: total}, nil
}

func (s *SubscriptionStoragePG) GetSubscriptionByID(ctx context.Context, id int) (*models.SubscriptionForStorage, error) {
	query := `
		SELECT service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`

	var res models.SubscriptionForStorage
	startDate := pgtype.Date{}
	endDate := pgtype.Date{}
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&res.ServiceName,
		&res.Price,
		&res.UserID,
		&startDate,
		&endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения подписки по id %d: %w", id, err)
	}

	dateString := startDate.Time
	res.StartDate = dateString

	if !endDate.Valid {
		res.EndDate = nil
	} else {
		dateString := endDate.Time
		res.EndDate = &dateString
	}

	return &res, nil
}

func (s *SubscriptionStoragePG) UpdateSubscriptionByID(ctx context.Context, req *models.SubscriptionForStorageWithId) error {
	var endDate pgtype.Date
	if req.EndDate != nil {
		endDate = pgtype.Date{Time: *req.EndDate, Valid: true}
	} else {
		endDate = pgtype.Date{Valid: false}
	}

	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, start_date = $3, end_date = $4
		WHERE id = $5
	`

	result, err := s.pool.Exec(ctx, query,
		req.ServiceName,
		req.Price,
		req.StartDate,
		endDate,
		req.ID,
	)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении подписки: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("подписка с id %d не найдена", req.ID)
	}

	return nil
}

func (s *SubscriptionStoragePG) DeleteSubscriptionByID(ctx context.Context, id int) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении подписки: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("подписка с id %d не найдена", id)
	}

	return nil
}

func (s *SubscriptionStoragePG) ListSubscriptionsByUserID(ctx context.Context, userID string) ([]*models.SubscriptionForStorageWithId, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
	`

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка подписок: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.SubscriptionForStorageWithId
	for rows.Next() {
		var sub models.SubscriptionForStorageWithId
		var startDate time.Time
		var endDate pgtype.Date

		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&startDate,
			&endDate,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		sub.StartDate = startDate
		if endDate.Valid {
			dateStr := endDate.Time
			sub.EndDate = &dateStr
		}

		subscriptions = append(subscriptions, &sub)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по результатам: %w", err)
	}

	return subscriptions, nil
}
