package read_repository

import (
	"context"
	"errors"
	"log"
	"time"
	"subscription-service/internal/db"
	"subscription-service/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type subscriptionReadRepo struct {
	db *gorm.DB
}

func NewSubscriptionReadRepo(db *gorm.DB) SubscriptionReadRepository {
	return &subscriptionReadRepo{db: db}
}

func (r *subscriptionReadRepo) GetByID(ctx context.Context, id uint) (*model.Subscription, error) {
	log.Printf("[SubscriptionRepo] GetByID called with id=%d", id)

	var sub model.Subscription
	err := r.db.WithContext(ctx).First(&sub, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("[SubscriptionRepo] GetByID id=%d not found", id)
			return nil, nil
		}
		log.Printf("[SubscriptionRepo] GetByID id=%d error: %v", id, err)
		return nil, err
	}

	log.Printf("[SubscriptionRepo] GetByID id=%d success", id)
	return &sub, nil
}

func (r *subscriptionReadRepo) GetByUserID(ctx context.Context, userID string) ([]model.Subscription, error) {
	log.Printf("[SubscriptionRepo] GetByUserID called with userID=%s", userID)

	uid, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("[SubscriptionRepo] GetByUserID userID=%s invalid UUID: %v", userID, err)
		return nil, err
	}

	var subs []model.Subscription
	err = r.db.WithContext(ctx).Where("user_id = ?", uid).Find(&subs).Error
	if err != nil {
		log.Printf("[SubscriptionRepo] GetByUserID userID=%s error: %v", userID, err)
		return nil, err
	}

	log.Printf("[SubscriptionRepo] GetByUserID userID=%s success, found %d subscriptions", userID, len(subs))
	return subs, nil
}

func (r *subscriptionReadRepo) GetAll(ctx context.Context) ([]model.Subscription, error) {
	log.Println("[SubscriptionRepo] GetAll called")

	var subs []model.Subscription
	err := r.db.WithContext(ctx).Find(&subs).Error
	if err != nil {
		log.Printf("[SubscriptionRepo] GetAll error: %v", err)
		return nil, err
	}

	log.Printf("[SubscriptionRepo] GetAll success, found %d subscriptions", len(subs))
	return subs, nil
}

func (r *subscriptionReadRepo) SumPriceByFilter(
	ctx context.Context,
	userID *string,
	serviceName *string,
	startDate, endDate time.Time,
) (float64, error) {
	log.Println("[SubscriptionRepo] SumPriceByFilter called")

	query := r.db.WithContext(ctx).Model(&model.Subscription{}).Select("SUM(price) as total_price")

	if userID != nil && *userID != "" {
		uid, err := uuid.Parse(*userID)
		if err != nil {
			log.Printf("[SubscriptionRepo] SumPriceByFilter invalid userID=%s: %v", *userID, err)
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
		log.Printf("[SubscriptionRepo] SumPriceByFilter error: %v", err)
		return 0, err
	}

	log.Printf("[SubscriptionRepo] SumPriceByFilter success, total_price=%.2f", total)
	return total, nil
}
