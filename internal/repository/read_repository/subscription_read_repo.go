package read_repository

import (
	"context"
	"time"
	"github.com/winnamu6/go-subscription-service/internal/model"
)

type SubscriptionReadRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Subscription, error)
	GetByUserID(ctx context.Context, userID string) ([]model.Subscription, error)
	GetAll(ctx context.Context) ([]model.Subscription, error)
	SumPriceByFilter(
		ctx context.Context,
		userID *string,
		serviceName *string,
		startDate time.Time,
		endDate time.Time,
	) (float64, error)
}
