package models

type SubscriptionWithId struct {
	ID int `json:"id"`
	Subscription
}
type Subscription struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type TotalCostRequest struct {
	UserID        *string  `json:"user_id,omitempty"`
	StartDate     string   `json:"start_date"`
	EndDate       *string  `json:"end_date"`
	Subscriptions []string `json:"subscriptions"`
}

type TotalCostResponse struct {
	Total int `json:"total"`
}
