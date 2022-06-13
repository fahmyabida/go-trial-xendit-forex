package model

import "time"

type (
	RateData struct {
		Success             bool      `json:"success"`
		Message             string    `json:"message"`
		OriginCurrency      string    `json:"origin_currency"`
		DestinationCurrency string    `json:"destination_currency"`
		Rate                float64   `json:"rate"`
		ExpiredAt           time.Time `json:"expired_at"`
	}
)

type (
	RequestBook struct {
		OriginCurrency      string `json:"origin_currency"`
		DestinationCurrency string `json:"destination_currency"`
	}
	ResponseBook struct {
		Success   bool      `json:"success"`
		Message   string    `json:"message"`
		Rate      float64   `json:"rate"`
		BookId    string    `json:"book_id"`
		ExpiredAt time.Time `json:"expired_at"`
	}
)
 