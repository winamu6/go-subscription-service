package read_repository

import (
	"context"
	"subscription-service/internal/model"
)

type SubscriptionReadRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Subscription, error)
	GetByUserID(ctx context.Context, userID string) ([]model.Subscription, error)
	GetAll(ctx context.Context) ([]model.Subscription, error)
}
