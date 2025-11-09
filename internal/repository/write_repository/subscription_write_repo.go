package write_repository

import (
	"context"
	"github.com/winnamu6/go-subscription-service/internal/model"
)

type SubscriptionWriteRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	Update(ctx context.Context, sub *model.Subscription) error
	Delete(ctx context.Context, id uint) error
}