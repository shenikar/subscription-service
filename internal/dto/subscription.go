package dto

import "github.com/google/uuid"

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required min=1"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required,datetime=01-2025"`
	EndDate     *string   `json:"end_date,omitempty" binding:"omitempty,datetime=01-2025"`
}

type UpdateSubcscriptionRequest struct {
	ServiceName *string    `json:"service_name" binding:"omitempty"`
	Price       *int       `json:"price" binding:"omitempty min=1"`
	UserID      *uuid.UUID `json:"user_id" binding:"omitempty"`
	StartDate   *string    `json:"start_date" binding:"omitempty,datetime=01-2025"`
	EndDate     *string    `json:"end_date,omitempty" binding:"omitempty,datetime=01-2025"`
}

type SubscriptionResponse struct {
	ID          int64     `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}
