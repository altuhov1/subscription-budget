package storage

import (
	"budget/internal/models"
	"context"
)

type SubscriptionStorage interface {
	CreateSubscription(ctx context.Context, req *models.SubscriptionForStorage) (int, error)
	GetSubscriptionsWithParam(ctx context.Context, req *models.TotalCostRequestForStorage) (*models.TotalCostResponse, error)
	GetSubscriptionByID(ctx context.Context, id int) (*models.SubscriptionForStorage, error)
	UpdateSubscriptionByID(ctx context.Context, req *models.SubscriptionForStorageWithId) error
	DeleteSubscriptionByID(ctx context.Context, id int) error
	ListSubscriptionsByUserID(ctx context.Context, userID string) ([]*models.SubscriptionForStorageWithId, error)
}
