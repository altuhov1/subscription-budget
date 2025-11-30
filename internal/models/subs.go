package models

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubscriptionPeriod struct {
	ServiceName string `json:"service_name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type TotalCostRequest struct {
	UserID        *string              `json:"user_id,omitempty"`
	Subscriptions []SubscriptionPeriod `json:"subscriptions"`
}

type TotalCostResponse struct {
	Total  int      `json:"total"`
	Errors []string `json:"errors"`
}
