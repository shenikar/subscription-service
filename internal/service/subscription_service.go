package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shenikar/subscription-service/internal/dto"
	"github.com/shenikar/subscription-service/internal/logger"
	"github.com/shenikar/subscription-service/internal/mapper"
	"github.com/shenikar/subscription-service/internal/model"
	"github.com/shenikar/subscription-service/internal/repository"
	"github.com/sirupsen/logrus"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, req dto.CreateSubscriptionRequest) (model.Subscription, error) {
	log := logger.GetLogger()
	sub, err := mapper.ToModelSubscription(req)
	if err != nil {
		log.WithError(err).Warn("Create: invalid subscription data")
		return model.Subscription{}, err
	}

	err = s.repo.Create(ctx, &sub)
	if err != nil {
		log.WithError(err).Error("failed to create subscription in repository")
		return model.Subscription{}, fmt.Errorf("could not create subscription: %w", err)
	}
	log.WithFields(logrus.Fields{
		"id":           sub.ID,
		"service_name": sub.ServiceName,
		"user_id":      sub.UserID,
	}).Info("subscription created successfully")

	return sub, nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id int64) (*model.Subscription, error) {
	log := logger.GetLogger()
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.WithError(err).Errorf("failed to get subscription by ID: %d", id)
		return nil, fmt.Errorf("get by id failed: %w", err)
	}
	if sub == nil {
		log.Warnf("subscription not found: %d", id)
		return nil, nil
	}

	return sub, nil
}

func (s *SubscriptionService) GelAll(ctx context.Context) ([]model.Subscription, error) {
	log := logger.GetLogger()

	subs, err := s.repo.GelAll(ctx)
	if err != nil {
		log.WithError(err).Error("failed to get subscriptions")
		return nil, fmt.Errorf("subscriptions failed: %w", err)
	}

	var res []model.Subscription
	for _, s := range subs {
		res = append(res, *s)
	}

	return res, nil
}
func (s *SubscriptionService) Update(ctx context.Context, id int64, req dto.UpdateSubscriptionRequest) (model.Subscription, error) {
	log := logger.GetLogger()

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.WithError(err).Errorf("failed to get subscription for update: %d", id)
		return model.Subscription{}, fmt.Errorf("get for update failed: %w", err)
	}
	if current == nil {
		log.Warnf("subscription to update not found: %d", id)
		return model.Subscription{}, fmt.Errorf("subscription not found")
	}

	updated, err := mapper.ToModelSubscriptionFromUpdate(id, req, *current)
	if err != nil {
		log.WithError(err).Error("failed to map update request")
		return model.Subscription{}, fmt.Errorf("invalid input: %w", err)
	}

	if err := s.repo.Update(ctx, &updated); err != nil {
		log.WithError(err).Errorf("failed to update subscription: %d", id)
		return model.Subscription{}, fmt.Errorf("update failed: %w", err)
	}

	log.WithFields(logrus.Fields{
		"id":      updated.ID,
		"user_id": updated.UserID,
	}).Info("subscription updated")

	return updated, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id int64) error {
	log := logger.GetLogger()

	if err := s.repo.Delete(ctx, id); err != nil {
		log.WithError(err).Errorf("failed to delete subscription: %d", id)
		return fmt.Errorf("delete failed: %w", err)
	}

	log.WithField("id", id).Info("subscription deleted")
	return nil
}

func (s *SubscriptionService) TotalPrice(ctx context.Context, req dto.TotalPriceFilterDTO) (int, error) {
	log := logger.GetLogger()

	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.WithError(err).Errorf("invalid user_id format: %s", req.UserID)
		return 0, fmt.Errorf("invalid user_id: %w", err)
	}
	sum, err := s.repo.TotalSumSubscription(ctx, &userUUID, &req.ServiceName, req.FromDate, req.ToDate)
	if err != nil {
		log.WithError(err).Error("failed to calculate total subscription price")
		return 0, fmt.Errorf("calculate total failed: %w", err)
	}

	log.WithFields(logrus.Fields{
		"user_id":      req.UserID,
		"service_name": req.ServiceName,
		"from":         req.FromDate,
		"to":           req.ToDate,
		"sum":          sum,
	}).Info("calculated total subscription price")

	return sum, nil
}
