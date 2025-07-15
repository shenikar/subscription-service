package mapper

import (
	"fmt"
	"time"

	"github.com/shenikar/subscription-service/internal/dto"
	"github.com/shenikar/subscription-service/internal/model"
)

func parseMonthYear(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}

func formatMonthYear(t time.Time) string {
	return t.Format("01-2006")
}

func ToModelSubscription(dto dto.CreateSubscriptionRequest) (model.Subscription, error) {
	startDate, err := parseMonthYear(dto.StartDate)
	if err != nil {
		return model.Subscription{}, fmt.Errorf("invalid start_date format: %w", err)
	}

	var endDate *time.Time
	if dto.EndDate != nil {
		ed, err := parseMonthYear(*dto.EndDate)
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
		s := formatMonthYear(*sub.EndDate)
		endDateSrt = &s
	}
	return dto.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   formatMonthYear(sub.StartDate),
		EndDate:     endDateSrt,
	}
}
