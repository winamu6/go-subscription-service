package read_repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/winnamu6/go-subscription-service/internal/logger"
	"github.com/winnamu6/go-subscription-service/internal/model"
	"gorm.io/gorm"
)

type subscriptionReadRepo struct {
	db *gorm.DB
}

func NewSubscriptionReadRepo(db *gorm.DB) SubscriptionReadRepository {
	return &subscriptionReadRepo{db: db}
}

func (r *subscriptionReadRepo) GetByID(ctx context.Context, id uint) (*model.Subscription, error) {
	log := logger.Get()
	log.Infof("[SubscriptionReadRepo] GetByID called | id=%d", id)

	var sub model.Subscription
	err := r.db.WithContext(ctx).First(&sub, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warnf("[SubscriptionReadRepo] GetByID not found | id=%d", id)
			return nil, nil
		}
		log.Errorf("[SubscriptionReadRepo] GetByID error | id=%d err=%v", id, err)
		return nil, err
	}

	log.Infof("[SubscriptionReadRepo] GetByID success | id=%d", id)
	return &sub, nil
}

func (r *subscriptionReadRepo) GetByUserID(ctx context.Context, userID string) ([]model.Subscription, error) {
	log := logger.Get()
	log.Infof("[SubscriptionReadRepo] GetByUserID called | userID=%s", userID)

	uid, err := uuid.Parse(userID)
	if err != nil {
		log.Warnf("[SubscriptionReadRepo] GetByUserID invalid UUID | userID=%s err=%v", userID, err)
		return nil, err
	}

	var subs []model.Subscription
	err = r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&subs).Error
	if err != nil {
		log.Errorf("[SubscriptionReadRepo] GetByUserID error | userID=%s err=%v", userID, err)
		return nil, err
	}

	log.Infof("[SubscriptionReadRepo] GetByUserID success | userID=%s count=%d", userID, len(subs))
	return subs, nil
}

func (r *subscriptionReadRepo) GetAll(ctx context.Context) ([]model.Subscription, error) {
	log := logger.Get()
	log.Info("[SubscriptionReadRepo] GetAll called")

	var subs []model.Subscription
	err := r.db.WithContext(ctx).Find(&subs).Error
	if err != nil {
		log.Errorf("[SubscriptionReadRepo] GetAll error | err=%v", err)
		return nil, err
	}

	log.Infof("[SubscriptionReadRepo] GetAll success | count=%d", len(subs))
	return subs, nil
}

func (r *subscriptionReadRepo) SumPriceByFilter(
	ctx context.Context,
	userID *string,
	serviceName *string,
	startDate, endDate time.Time,
) (float64, error) {
	log := logger.Get()
	log.Infof("[SubscriptionReadRepo] SumPriceByFilter called | userID=%v serviceName=%v start=%s end=%s",
		userID, serviceName, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))

	query := r.db.WithContext(ctx).Model(&model.Subscription{}).Select("SUM(price) as total_price")

	if userID != nil && *userID != "" {
		uid, err := uuid.Parse(*userID)
		if err != nil {
			log.Warnf("[SubscriptionReadRepo] SumPriceByFilter invalid userID=%s err=%v", *userID, err)
			return 0, err
		}
		query = query.Where("user_id = ?", uid)
	}

	if serviceName != nil && *serviceName != "" {
		query = query.Where("service_name = ?", *serviceName)
	}

	query = query.Where("start_date >= ? AND start_date <= ?", startDate, endDate)

	var total float64
	err := query.Scan(&total).Error
	if err != nil {
		log.Errorf("[SubscriptionReadRepo] SumPriceByFilter error | err=%v", err)
		return 0, err
	}

	log.Infof("[SubscriptionReadRepo] SumPriceByFilter success | total_price=%.2f", total)
	return total, nil
}