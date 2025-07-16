package mapper

import (
	"fmt"
	"time"

	"github.com/shenikar/subscription-service/internal/dto"
	"github.com/shenikar/subscription-service/internal/model"
)

func ParseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}

func FormatMonthYear(t time.Time) string {
	return t.Format("01-2006")
}

func ToModelSubscription(dto dto.CreateSubscriptionRequest) (model.Subscription, error) {
	startDate, err := ParseMonthYear(dto.StartDate)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("invalid start_date format: %w", err)
	}

	var endDate *time.Time
	if dto.EndDate != nil {
		ed, err := ParseMonthYear(*dto.EndDate)
		if err != nil {
			return model.Subscription{}, fmt.Errorf("invalid end_date format: %w", err)
		}
		endDate = &ed
	}
	return model.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      dto.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func ToResponseDTO(sub model.Subscription) dto.SubscriptionResponse {
	var endDateSrt *string
	if sub.EndDate != nil {
		s := FormatMonthYear(*sub.EndDate)
		endDateSrt = &s
	}
	return dto.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   FormatMonthYear(sub.StartDate),
		EndDate:     endDateSrt,
	}
}

func ToModelSubscriptionFromUpdate(id int64, dto dto.UpdateSubscriptionRequest, current model.Subscription) (model.Subscription, error) {
	sub := current
	sub.ID = id

	if dto.ServiceName != nil {
		sub.ServiceName = *dto.ServiceName
	}
	if dto.Price != nil {
		sub.Price = *dto.Price
	}
	if dto.UserID != nil {
		sub.UserID = *dto.UserID
	}
	if dto.StartDate != nil {
		startDate, err := ParseMonthYear(*dto.StartDate)
		if err != nil {
			return model.Subscription{}, fmt.Errorf("invalid start_date format: %w", err)
		}
		sub.StartDate = startDate
	}
	if dto.EndDate != nil {
		if *dto.EndDate != "" {
			endDate, err := ParseMonthYear(*dto.EndDate)
			if err != nil {
				return model.Subscription{}, fmt.Errorf("invalid end_date format: %w", err)
			}
			sub.EndDate = &endDate
		} else {
			sub.EndDate = nil
		}
	}

	return sub, nil
}
