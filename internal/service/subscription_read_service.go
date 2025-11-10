package service

import (
	"context"
	"time"

	"github.com/winnamu6/go-subscription-service/internal/logger"
	"github.com/winnamu6/go-subscription-service/internal/model"
	"github.com/winnamu6/go-subscription-service/internal/repository/read_repository"
)

type SubscriptionQueryService interface {
	GetByID(ctx context.Context, id uint) (*model.SubscriptionResponse, error)
	GetByUserID(ctx context.Context, userID string) ([]model.SubscriptionResponse, error)
	GetAll(ctx context.Context) ([]model.SubscriptionResponse, error)
	SumPriceByFilter(ctx context.Context, userID *string, serviceName *string, startDate, endDate time.Time) (int, error)
}

type subscriptionQueryService struct {
	readRepo read_repository.SubscriptionReadRepository
}

func NewSubscriptionQueryService(readRepo read_repository.SubscriptionReadRepository) SubscriptionQueryService {
	return &subscriptionQueryService{readRepo: readRepo}
}

func (s *subscriptionQueryService) GetByID(ctx context.Context, id uint) (*model.SubscriptionResponse, error) {
	log := logger.Get()
	log.Infof("[QueryService] GetByID called | id=%d", id)

	sub, err := s.readRepo.GetByID(ctx, id)
	if err != nil {
		log.Errorf("[QueryService] GetByID error | id=%d err=%v", id, err)
		return nil, err
	}

	log.Infof("[QueryService] GetByID success | id=%d", id)
	return toSubscriptionResponse(sub), nil
}

func (s *subscriptionQueryService) GetByUserID(ctx context.Context, userID string) ([]model.SubscriptionResponse, error) {
	log := logger.Get()
	log.Infof("[QueryService] GetByUserID called | userID=%s", userID)

	subs, err := s.readRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Errorf("[QueryService] GetByUserID error | userID=%s err=%v", userID, err)
		return nil, err
	}

	log.Infof("[QueryService] GetByUserID success | userID=%s count=%d", userID, len(subs))
	return toSubscriptionResponseList(subs), nil
}

func (s *subscriptionQueryService) GetAll(ctx context.Context) ([]model.SubscriptionResponse, error) {
	log := logger.Get()
	log.Info("[QueryService] GetAll called")

	subs, err := s.readRepo.GetAll(ctx)
	if err != nil {
		log.Errorf("[QueryService] GetAll error | err=%v", err)
		return nil, err
	}

	log.Infof("[QueryService] GetAll success | count=%d", len(subs))
	return toSubscriptionResponseList(subs), nil
}

func (s *subscriptionQueryService) SumPriceByFilter(ctx context.Context, userID *string, serviceName *string, startDate, endDate time.Time) (int, error) {
	log := logger.Get()
	log.Infof("[QueryService] SumPriceByFilter called | userID=%v serviceName=%v start=%s end=%s",
		userID, serviceName, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))

	total, err := s.readRepo.SumPriceByFilter(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		log.Errorf("[QueryService] SumPriceByFilter error | err=%v", err)
		return 0, err
	}

	log.Infof("[QueryService] SumPriceByFilter success | total=%d", int(total))
	return int(total), nil
}

func toSubscriptionResponse(sub *model.Subscription) *model.SubscriptionResponse {
	if sub == nil {
		return nil
	}
	return &model.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:		 float64(sub.Price),
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}
}

func toSubscriptionResponseList(subs []model.Subscription) []model.SubscriptionResponse {
	res := make([]model.SubscriptionResponse, len(subs))
	for i, sub := range subs {
		res[i] = *toSubscriptionResponse(&sub)
	}
	return res
}
