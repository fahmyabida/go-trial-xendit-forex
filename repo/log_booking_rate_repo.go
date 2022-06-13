package repo

import (
	"fmt"
	"go-forex/model"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

type LogBookingRateRepo struct {
	db *pg.DB
}

func NewLogBookingRateRepo(db *pg.DB) LogBookingRateRepo {
	return LogBookingRateRepo{db}
}

func (r LogBookingRateRepo) InsertLogBookingRateRepo(traceId, bookId, origin, destination string, rate float64,
	expireAt time.Time) error {
	data := model.LogBookingRate{
		TraceId:             traceId,
		BookId:              bookId,
		OriginCurrency:      origin,
		DestinationCurrency: destination,
		Rate:                rate,
		ExpiredAt:           expireAt,
		CreatedAt:           time.Time{},
	}
	if _, err := r.db.Model(&data).Insert(); err != nil {
		log.Err(fmt.Errorf("%v, %v", traceId, err))
		return err
	}
	return nil
}
