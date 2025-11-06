package model

import (
	"time"

	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required,min=2,max=255"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"omitempty,min=2,max=255"`
	Price       *float64   `json:"price,omitempty" binding:"omitempty,gt=0"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubscriptionResponse struct {
	ID          uint       `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       float64    `json:"price"`
	UserID      uuid.UUID  `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
