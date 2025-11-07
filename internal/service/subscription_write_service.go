package service

import (
	"context"
	"errors"
	"subscription-service/internal/model"
	"subscription-service/internal/repository/write_repository"
	"time"
)

type SubscriptionCommandService interface {
	Create(ctx context.Context, req *model.CreateSubscriptionRequest) (*model.SubscriptionResponse, error)
	Update(ctx context.Context, id uint, req *model.UpdateSubscriptionRequest) (*model.SubscriptionResponse, error)
	Delete(ctx context.Context, id uint) error
}

type subscriptionCommandService struct {
	writeRepo write_repository.SubscriptionWriteRepository
	readSvc   SubscriptionQueryService
}

func NewSubscriptionCommandService(writeRepo write_repository.SubscriptionWriteRepository, readSvc SubscriptionQueryService) SubscriptionCommandService {
	return &subscriptionCommandService{
		writeRepo: writeRepo,
		readSvc:   readSvc,
	}
}

func (s *subscriptionCommandService) Create(ctx context.Context, req *model.CreateSubscriptionRequest) (*model.SubscriptionResponse, error) {
	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return nil, errors.New("end_date cannot be before start_date")
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.writeRepo.Create(ctx, sub); err != nil {
		return nil, err
	}

	return toSubscriptionResponse(sub), nil
}

func (s *subscriptionCommandService) Update(ctx context.Context, id uint, req *model.UpdateSubscriptionRequest) (*model.SubscriptionResponse, error) {
	existing, err := s.readSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("subscription not found")
	}

	sub := &model.Subscription{
		ID:          id,
		ServiceName: coalesceString(req.ServiceName, existing.ServiceName),
		Price:       coalesceFloat(req.Price, existing.Price),
		UserID:      existing.UserID,
		StartDate:   coalesceTime(req.StartDate, existing.StartDate),
		EndDate:     coalesceTimePtr(req.EndDate, existing.EndDate),
		UpdatedAt:   time.Now(),
	}

	if err := s.writeRepo.Update(ctx, sub); err != nil {
		return nil, err
	}

	return toSubscriptionResponse(sub), nil
}

func (s *subscriptionCommandService) Delete(ctx context.Context, id uint) error {
	existing, err := s.readSvc.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("subscription not found")
	}

	return s.writeRepo.Delete(ctx, id)
}


func coalesceString(newVal, oldVal string) string {
	if newVal != "" {
		return newVal
	}
	return oldVal
}

func coalesceFloat(newVal *float64, oldVal float64) float64 {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func coalesceTime(newVal *time.Time, oldVal time.Time) time.Time {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func coalesceTimePtr(newVal, oldVal *time.Time) *time.Time {
	if newVal != nil {
		return newVal
	}
	return oldVal
}
