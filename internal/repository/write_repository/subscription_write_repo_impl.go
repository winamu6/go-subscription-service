package write_repository

import (
	"context"
	"log"
	"time"
	"subscription-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type subscriptionWriteRepo struct {
	db *gorm.DB
}

func NewSubscriptionWriteRepo(db *gorm.DB) SubscriptionWriteRepository {
	return &subscriptionWriteRepo{db: db}
}

func (r *subscriptionWriteRepo) Create(ctx context.Context, sub *model.Subscription) error {
	log.Printf("[SubscriptionWriteRepo] Create called for userID=%s, serviceName=%s", sub.UserID, sub.ServiceName)

	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()

	if sub.ID != 0 {
		log.Printf("[SubscriptionWriteRepo] Create warning: ID should be zero for new subscription, got %d", sub.ID)
	}

	if err := r.db.WithContext(ctx).Create(sub).Error; err != nil {
		log.Printf("[SubscriptionWriteRepo] Create error: %v", err)
		return err
	}

	log.Printf("[SubscriptionWriteRepo] Create success: subscription ID=%d", sub.ID)
	return nil
}


func (r *subscriptionWriteRepo) Update(ctx context.Context, sub *model.Subscription) error {
	log.Printf("[SubscriptionWriteRepo] Update called for subscription ID=%d", sub.ID)

	sub.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(sub).Error; err != nil {
		log.Printf("[SubscriptionWriteRepo] Update error: %v", err)
		return err
	}

	log.Printf("[SubscriptionWriteRepo] Update success for subscription ID=%d", sub.ID)
	return nil
}


func (r *subscriptionWriteRepo) Delete(ctx context.Context, id uint) error {
	log.Printf("[SubscriptionWriteRepo] Delete called for subscription ID=%d", id)

	if err := r.db.WithContext(ctx).Delete(&model.Subscription{}, id).Error; err != nil {
		log.Printf("[SubscriptionWriteRepo] Delete error: %v", err)
		return err
	}

	log.Printf("[SubscriptionWriteRepo] Delete success for subscription ID=%d", id)
	return nil
}