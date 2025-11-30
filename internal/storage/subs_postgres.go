package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"budget/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamPostgresStorage struct {
	pool *pgxpool.Pool
}

func NewSubsStorage(pool *pgxpool.Pool) *TeamPostgresStorage {
	return &TeamPostgresStorage{pool: pool}
}

// CreateSubscription сохраняет новую подписку в БД
func (s *TeamPostgresStorage) CreateSubscription(ctx context.Context, req models.CreateSubscriptionRequest) error {
	// Преобразуем start_date из MM-YYYY в дату (первое число месяца)
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date format: %w", err)
	}

	var endDate *time.Time
	if req.EndDate != nil {
		// Преобразуем end_date в последнее число месяца
		end, err := parseMonthYearEnd(*req.EndDate)
		if err != nil {
			return fmt.Errorf("invalid end_date format: %w", err)
		}
		endDate = &end
	}

	sql := `
	INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err = s.pool.Exec(ctx, sql,
		req.ServiceName,
		req.Price,
		req.UserID,
		startDate,
		endDate,
	)
	return err
}

// GetTotalCost считает суммарную стоимость подписок с детализацией ошибок
func (s *TeamPostgresStorage) GetTotalCost(ctx context.Context, req models.TotalCostRequest) (models.TotalCostResponse, error) {
	var total int
	var errorsMas []string

	// Обрабатываем каждую подписку из запроса
	for _, sub := range req.Subscriptions {
		// Парсим даты периода
		periodStart, err := time.Parse("01-2006", sub.StartDate)
		if err != nil {
			errorsMas = append(errorsMas, fmt.Sprintf("invalid start_date %s for service %s: %v", sub.StartDate, sub.ServiceName, err))
			continue
		}

		periodEnd, err := parseMonthYearEnd(sub.EndDate)
		if err != nil {
			errorsMas = append(errorsMas, fmt.Sprintf("invalid end_date %s for service %s: %v", sub.EndDate, sub.ServiceName, err))
			continue
		}

		// Формируем условия фильтрации
		params := []interface{}{sub.ServiceName, periodEnd, periodStart}
		query := `
		SELECT COALESCE(SUM(price), 0), COUNT(*)
		FROM subscriptions
		WHERE service_name = $1
		AND start_date <= $2
		AND (end_date IS NULL OR end_date >= $3)
		`

		// Добавляем фильтр по пользователю если указан
		if req.UserID != nil {
			params = append(params, *req.UserID)
			query += ` AND user_id = $` + fmt.Sprint(len(params))
		}

		var sum int
		var count int
		err = s.pool.QueryRow(ctx, query, params...).Scan(&sum, &count)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				errorsMas = append(errorsMas, fmt.Sprintf("database returned no rows for service %s", sub.ServiceName))
			} else {
				return models.TotalCostResponse{}, fmt.Errorf("database query failed: %w", err)
			}
		}

		if count == 0 {
			errorsMas = append(errorsMas, fmt.Sprintf(
				"no active subscriptions found for service '%s' in period %s - %s",
				sub.ServiceName,
				sub.StartDate,
				sub.EndDate,
			))
		}

		total += sum
	}

	return models.TotalCostResponse{
		Total:  total,
		Errors: errorsMas,
	}, nil
}

// Вспомогательная функция для получения последнего дня месяца
func parseMonthYearEnd(input string) (time.Time, error) {
	t, err := time.Parse("01-2006", input)
	if err != nil {
		return time.Time{}, err
	}
	// Получаем последний день следующего месяца и вычитаем один день
	return time.Date(t.Year(), t.Month()+1, 0, 23, 59, 59, 999999999, time.UTC), nil
}
