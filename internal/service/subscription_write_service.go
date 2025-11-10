package service

import (
	"context"
	"errors"
	"time"

	"github.com/winnamu6/go-subscription-service/internal/logger"
	"github.com/winnamu6/go-subscription-service/internal/model"
	"github.com/winnamu6/go-subscription-service/internal/repository/write_repository"
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
	log := logger.Get()
	log.Infof("[CommandService] Create called | userID=%s serviceName=%s", req.UserID, req.ServiceName)

	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		log.Warnf("[CommandService] Create invalid dates | start=%s end=%s", req.StartDate, req.EndDate)
		return nil, errors.New("end_date cannot be before start_date")
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       int(req.Price),
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := s.writeRepo.Create(ctx, sub); err != nil {
		log.Errorf("[CommandService] Create error | userID=%s err=%v", req.UserID, err)
		return nil, err
	}

	log.Infof("[CommandService] Create success | subscriptionID=%d", sub.ID)
	return toSubscriptionResponse(sub), nil
}

func (s *subscriptionCommandService) Update(ctx context.Context, id uint, req *model.UpdateSubscriptionRequest) (*model.SubscriptionResponse, error) {
	log := logger.Get()
	log.Infof("[CommandService] Update called | id=%d", id)

	existing, err := s.readSvc.GetByID(ctx, id)
	if err != nil {
		log.Errorf("[CommandService] Update read error | id=%d err=%v", id, err)
		return nil, err
	}
	if existing == nil {
		log.Warnf("[CommandService] Update failed | id=%d not found", id)
		return nil, errors.New("subscription not found")
	}

	startDate := coalesceTime(req.StartDate, existing.StartDate)
	endDate := coalesceTimePtr(req.EndDate, existing.EndDate)
	if endDate != nil && endDate.Before(startDate) {
		log.Warnf("[CommandService] Update invalid dates | id=%d start=%s end=%s", id, startDate, *endDate)
		return nil, errors.New("end_date cannot be before start_date")
	}

	price := coalesceFloatToInt(req.Price, int(existing.Price))

	sub := &model.Subscription{
		ID:          id,
		ServiceName: coalesceString(req.ServiceName, existing.ServiceName),
		Price:       price,
		UserID:      existing.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   existing.CreatedAt,
	}

	if err := s.writeRepo.Update(ctx, sub); err != nil {
		log.Errorf("[CommandService] Update error | id=%d err=%v", id, err)
		return nil, err
	}

	log.Infof("[CommandService] Update success | id=%d", id)
	return toSubscriptionResponse(sub), nil
}

func (s *subscriptionCommandService) Delete(ctx context.Context, id uint) error {
	log := logger.Get()
	log.Infof("[CommandService] Delete called | id=%d", id)

	existing, err := s.readSvc.GetByID(ctx, id)
	if err != nil {
		log.Errorf("[CommandService] Delete read error | id=%d err=%v", id, err)
		return err
	}
	if existing == nil {
		log.Warnf("[CommandService] Delete failed | id=%d not found", id)
		return errors.New("subscription not found")
	}

	if err := s.writeRepo.Delete(ctx, id); err != nil {
		log.Errorf("[CommandService] Delete error | id=%d err=%v", id, err)
		return err
	}

	log.Infof("[CommandService] Delete success | id=%d", id)
	return nil
}

func coalesceString(newVal, oldVal string) string {
	if newVal != "" {
		return newVal
	}
	return oldVal
}

func coalesceInt(newVal *int, oldVal int) int {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func coalesceFloatToInt(newVal *float64, oldVal int) int {
	if newVal != nil {
		return int(*newVal)
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