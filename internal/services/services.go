package services

import (
	"budget/internal/models"
	"context"
)

type SubManager interface {
	CreateSubscription(ctx context.Context, req *models.Subscription) (int, error)
	GetSubscriptionsWithParam(ctx context.Context, req *models.TotalCostRequest) (*models.TotalCostResponse, error)
	GetSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubscriptionByID(ctx context.Context, req *models.SubscriptionWithId) error
	DeleteSubscriptionByID(ctx context.Context, id int) error
	ListSubscriptionsByUserID(ctx context.Context, userID string) ([]*models.SubscriptionWithId, error)
}
