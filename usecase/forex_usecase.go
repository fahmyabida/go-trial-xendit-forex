package usecase

import (
	http_client "go-forex/http/client"
	"go-forex/model"
	"go-forex/redis_jobs"
	"go-forex/repo"
	"time"

	"github.com/google/uuid"
)

type GetRateUsecase struct {
	httpClient         http_client.HttpClient
	redisClient        redis_jobs.RedisClient
	logBookingRateRepo repo.LogBookingRateRepo
}

func NewGetRateUsecase(httpClient http_client.HttpClient,
	redisClient redis_jobs.RedisClient,
	logBookingRateRepo repo.LogBookingRateRepo) GetRateUsecase {
	return GetRateUsecase{httpClient,
		redisClient,
		logBookingRateRepo}
}

func (u GetRateUsecase) GetRate(traceId, origin, destination, tokenAccess string) (rateData model.RateData) {
	rateData = u.httpClient.GetForex(traceId, origin, destination)
	if !rateData.Success {
		return rateData
	}
	rateData.Message = "success"
	return rateData
}

func (u GetRateUsecase) LockOrBookRate(traceId, tokenAccess string, requestBook model.RequestBook) (response model.ResponseBook) {
	response.Rate, response.ExpiredAt = u.redisClient.GetCurrentRate(traceId, requestBook.OriginCurrency, requestBook.DestinationCurrency)
	if response.Rate <= 0 {
		response.Message = "rate are expired"
		return response
	}
	response.BookId = uuid.New().String()
	ThirtyMinuteAhead := time.Now().Add(30 * time.Minute)
	if ThirtyMinuteAhead.Before(response.ExpiredAt) {
		response.ExpiredAt = ThirtyMinuteAhead
	}
	err := u.redisClient.SetBookRate(traceId, response.BookId, requestBook.OriginCurrency, requestBook.DestinationCurrency, response.Rate,
		response.ExpiredAt)
	if err != nil {
		response.Message = "unexpected error on memcache"
		return response
	}
	err = u.logBookingRateRepo.InsertLogBookingRateRepo(traceId, response.BookId, requestBook.OriginCurrency, requestBook.DestinationCurrency, response.Rate,
		response.ExpiredAt)
	if err != nil {
		response.Message = "unexpected error on database"
		return response
	}
	response.Success = true
	response.Message = "success"
	return response
}
