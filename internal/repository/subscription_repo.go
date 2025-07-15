package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/shenikar/subscription-service/internal/model"
)

type SubscriptionRepository struct {
	conn *pgx.Conn
}

func NewSubscriptionRepository(conn *pgx.Conn) *SubscriptionRepository {
	return &SubscriptionRepository{conn: conn}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *model.Subscription) error {
	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id;		
	`
	var id int64
	err := r.conn.QueryRow(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed insert subscription: %w", err)
	}
	sub.ID = id
	return nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id int64) (*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`

	var sub model.Subscription
	err := r.conn.QueryRow(ctx, query, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return &sub, nil
}
