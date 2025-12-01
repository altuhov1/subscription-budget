package models

import "time"

type SubscriptionForStorage struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
}

type SubscriptionForStorageWithId struct {
	ID int
	SubscriptionForStorage
}

type TotalCostRequestForStorage struct {
	UserID        *string
	StartDate     time.Time
	EndDate       *time.Time
	Subscriptions []string
}
