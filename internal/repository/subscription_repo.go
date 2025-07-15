package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func (r *SubscriptionRepository) GelAll(ctx context.Context) ([]*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscription: %w", err)
	}
	defer rows.Close()

	var subs []*model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subs = append(subs, &sub)
	}
	return subs, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *model.Subscription) error {
	query := `UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6	
	`
	_, err := r.conn.Exec(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	_, err := r.conn.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionRepository) TotalSumSubscription(ctx context.Context, userID *uuid.UUID, serviceName *string, from, to time.Time) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0)
		FROM subscriptions WHERE start_date >= $1 AND start_date <= $2
	`

	args := []interface{}{from, to}
	argNum := 3
	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, *userID)
		argNum++
	}

	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name ILIKE $%d", argNum)
		args = append(args, "%"+*serviceName+"%")
		argNum++
	}

	var sum int
	err := r.conn.QueryRow(ctx, query, args...).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total: %w", err)
	}
	return sum, nil
}
