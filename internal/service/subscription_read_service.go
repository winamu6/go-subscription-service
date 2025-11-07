package service

import (
	"context"
	"subscription-service/internal/model"
	"subscription-service/internal/repository/read_repository"
	"time"
)

type SubscriptionQueryService interface {
	GetByID(ctx context.Context, id uint) (*model.SubscriptionResponse, error)
	GetByUserID(ctx context.Context, userID string) ([]model.SubscriptionResponse, error)
	GetAll(ctx context.Context) ([]model.SubscriptionResponse, error)
	SumPriceByFilter(ctx context.Context, userID *string, serviceName *string, startDate, endDate time.Time) (float64, error)
}

type subscriptionQueryService struct {
	readRepo read_repository.SubscriptionReadRepository
}

func NewSubscriptionQueryService(readRepo read_repository.SubscriptionReadRepository) SubscriptionQueryService {
	return &subscriptionQueryService{readRepo: readRepo}
}

func (s *subscriptionQueryService) GetByID(ctx context.Context, id uint) (*model.SubscriptionResponse, error) {
	sub, err := s.readRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toSubscriptionResponse(sub), nil
}

func (s *subscriptionQueryService) GetByUserID(ctx context.Context, userID string) ([]model.SubscriptionResponse, error) {
	subs, err := s.readRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toSubscriptionResponseList(subs), nil
}

func (s *subscriptionQueryService) GetAll(ctx context.Context) ([]model.SubscriptionResponse, error) {
	subs, err := s.readRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return toSubscriptionResponseList(subs), nil
}

func (s *subscriptionQueryService) SumPriceByFilter(ctx context.Context, userID *string, serviceName *string, startDate, endDate time.Time) (float64, error) {
	return s.readRepo.SumPriceByFilter(ctx, userID, serviceName, startDate, endDate)
}


func toSubscriptionResponse(sub *model.Subscription) *model.SubscriptionResponse {
	if sub == nil {
		return nil
	}

	return &model.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}
}

func toSubscriptionResponseList(subs []model.Subscription) []model.SubscriptionResponse {
	res := make([]model.SubscriptionResponse, 0, len(subs))
	for _, sub := range subs {
		res = append(res, *toSubscriptionResponse(&sub))
	}
	return res
}