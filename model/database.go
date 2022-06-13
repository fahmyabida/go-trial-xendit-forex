package model

import "time"

type LogBookingRate struct {
	tableName           struct{}  `pg:"public.log_booking_rate"`
	TraceId             string    `pg:"trace_id"`
	BookId              string    `pg:"book_id"`
	OriginCurrency      string    `pg:"origin_currency"`
	DestinationCurrency string    `pg:"destination_currency"`
	Rate                float64   `pg:"rate"`
	ExpiredAt           time.Time `pg:"expired_at"`
	CreatedAt           time.Time `pg:"created_at,default:now()"`
	// ClientId            string    `pg:"client_id"`
}
