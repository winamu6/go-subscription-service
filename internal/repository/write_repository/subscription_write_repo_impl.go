package write_repository

import (
	"context"
	"time"

	"github.com/winnamu6/go-subscription-service/internal/logger"
	"github.com/winnamu6/go-subscription-service/internal/model"
	"gorm.io/gorm"
)

type subscriptionWriteRepo struct {
	db *gorm.DB
}

func NewSubscriptionWriteRepo(db *gorm.DB) SubscriptionWriteRepository {
	return &subscriptionWriteRepo{db: db}
}

func (r *subscriptionWriteRepo) Create(ctx context.Context, sub *model.Subscription) error {
	log := logger.Get()
	log.Infof("[SubscriptionWriteRepo] Create called | userID=%s serviceName=%s", sub.UserID, sub.ServiceName)

	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()

	if sub.ID != 0 {
		log.Warnf("[SubscriptionWriteRepo] Create warning: ID should be zero for new subscription (got %d)", sub.ID)
	}

	if err := r.db.WithContext(ctx).Create(sub).Error; err != nil {
		log.Errorf("[SubscriptionWriteRepo] Create error | userID=%s serviceName=%s err=%v", sub.UserID, sub.ServiceName, err)
		return err
	}

	log.Infof("[SubscriptionWriteRepo] Create success | subscriptionID=%d", sub.ID)
	return nil
}

func (r *subscriptionWriteRepo) Update(ctx context.Context, sub *model.Subscription) error {
	log := logger.Get()
	log.Infof("[SubscriptionWriteRepo] Update called | subscriptionID=%d", sub.ID)

	sub.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(sub).Error; err != nil {
		log.Errorf("[SubscriptionWriteRepo] Update error | subscriptionID=%d err=%v", sub.ID, err)
		return err
	}

	log.Infof("[SubscriptionWriteRepo] Update success | subscriptionID=%d", sub.ID)
	return nil
}

func (r *subscriptionWriteRepo) Delete(ctx context.Context, id uint) error {
	log := logger.Get()
	log.Infof("[SubscriptionWriteRepo] Delete called | subscriptionID=%d", id)

	if err := r.db.WithContext(ctx).Delete(&model.Subscription{}, id).Error; err != nil {
		log.Errorf("[SubscriptionWriteRepo] Delete error | subscriptionID=%d err=%v", id, err)
		return err
	}

	log.Infof("[SubscriptionWriteRepo] Delete success | subscriptionID=%d", id)
	return nil
}