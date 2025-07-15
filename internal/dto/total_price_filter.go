package dto

import (
	"time"
)

type TotalPriceFilterDTO struct {
	UserID      string    `form:"user_id"`
	ServiceName string    `form:"service_name"`
	FromDate    time.Time `form:"from_date" time_format:"02-01-2006"`
	ToDate      time.Time `form:"to_date" time_format:"02-01-2006"`
}
