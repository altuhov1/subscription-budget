package services

import (
	"budget/internal/models"
	"budget/internal/storage"
	"context"
	"fmt"
	"time"
)

type SubService struct {
	SubStorage storage.SubscriptionStorage
}

func NewPullRequestService(SubStorage storage.SubscriptionStorage) *SubService {
	return &SubService{
		SubStorage: SubStorage,
	}
}

func (s *SubService) CreateSubscription(ctx context.Context, req *models.Subscription) (int, error) {
	startDateParsed, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return -1, fmt.Errorf("некорректный формат start_date (ожидается MM-YYYY): %w", err)
	}
	var t time.Time
	var endDateParsed *time.Time
	if req.EndDate != nil {
		t, err = time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return -1, fmt.Errorf("некорректный формат end_date (ожидается MM-YYYY): %w", err)
		}
		endDateParsed = &t
	} else {
		endDateParsed = nil
	}
	reqS := &models.SubscriptionForStorage{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDateParsed,
		EndDate:     endDateParsed,
	}
	return s.SubStorage.CreateSubscription(ctx, reqS)
}

func (s *SubService) GetSubscriptionsWithParam(ctx context.Context, req *models.TotalCostRequest) (*models.TotalCostResponse, error) {
	startDateParsed, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("некорректный формат start_date (ожидается MM-YYYY): %w", err)
	}
	var t time.Time
	var endDateParsed *time.Time
	if req.EndDate != nil {
		t, err = time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("некорректный формат end_date (ожидается MM-YYYY): %w", err)
		}
		endDateParsed = &t
	} else {
		now := time.Now()
		t := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDateParsed = &t
	}
	reqS := &models.TotalCostRequestForStorage{
		UserID:    req.UserID,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}
	return s.SubStorage.GetSubscriptionsWithParam(ctx, reqS)
}
func (s *SubService) GetSubscriptionByID(ctx context.Context, id int) (*models.Subscription, error) {
	resS, err := s.SubStorage.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	var EndTime *string
	if resS.EndDate == nil {
		EndTime = nil
	} else {
		t := (*resS.EndDate).Format("01-2006")
		EndTime = &t
	}

	res := &models.Subscription{
		ServiceName: resS.ServiceName,
		Price:       resS.Price,
		UserID:      resS.UserID,
		StartDate:   resS.StartDate.Format("01-2006"),
		EndDate:     EndTime,
	}
	return res, nil
}
func (s *SubService) UpdateSubscriptionByID(ctx context.Context, req *models.SubscriptionWithId) error {
	startDateParsed, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return fmt.Errorf("некорректный формат start_date (ожидается MM-YYYY): %w", err)
	}
	var t time.Time
	var endDateParsed *time.Time
	if req.EndDate != nil {
		t, err = time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return fmt.Errorf("некорректный формат end_date (ожидается MM-YYYY): %w", err)
		}
		endDateParsed = &t
	} else {
		endDateParsed = nil
	}
	reqS := &models.SubscriptionForStorageWithId{
		ID: req.ID,
		SubscriptionForStorage: models.SubscriptionForStorage{
			ServiceName: req.ServiceName,
			Price:       req.Price,
			UserID:      req.UserID,
			StartDate:   startDateParsed,
			EndDate:     endDateParsed,
		},
	}
	return s.SubStorage.UpdateSubscriptionByID(ctx, reqS)
}
func (s *SubService) DeleteSubscriptionByID(ctx context.Context, id int) error {
	return s.SubStorage.DeleteSubscriptionByID(ctx, id)
}
func (s *SubService) ListSubscriptionsByUserID(ctx context.Context, userID string) ([]*models.SubscriptionWithId, error) {
	resS, err := s.SubStorage.ListSubscriptionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	res := make([]*models.SubscriptionWithId, 0, len(resS))
	for i := 0; i < len(resS); i++ {
		var EndTime *string
		if resS[i].EndDate == nil {
			EndTime = nil
		} else {
			t := (resS[i].EndDate).Format("01-2006")
			EndTime = &t
		}
		tempNode := &models.SubscriptionWithId{
			ID: resS[i].ID,
			Subscription: models.Subscription{
				ServiceName: resS[i].ServiceName,
				Price:       resS[i].Price,
				UserID:      resS[i].UserID,
				StartDate:   resS[i].StartDate.Format("01-2006"),
				EndDate:     EndTime,
			},
		}
		res = append(res, tempNode)
	}
	return res, nil

}
