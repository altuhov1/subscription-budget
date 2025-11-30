package storage

import (
	"budget/internal/models"
	"context"
)

type SubscriptionStorage interface {
	CreateSubscription(ctx context.Context, req *models.CreateSubscriptionRequest) (int, error)
	GetSubscriptionsWithParam(ctx context.Context, req *models.TotalCostRequest) (*models.TotalCostResponse, error)
	GetSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubscriptionByID(ctx context.Context, id int, req *models.CreateSubscriptionRequest) error
	DeleteSubscriptionByID(ctx context.Context, id int) error
	ListSubscriptionsByUserID(ctx context.Context, userID string) ([]*models.Subscription, error)
}
