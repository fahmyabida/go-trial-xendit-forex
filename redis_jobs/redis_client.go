package redis_jobs

import (
	"fmt"
	"go-forex/model"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog/log"
)

type RedisClient struct {
	redisClient *redis.Client
	svcConfig   model.SvcConfig
}

func NewRedisClient(redisClient *redis.Client, svcConfig model.SvcConfig) RedisClient {
	return RedisClient{redisClient, svcConfig}
}

func (r RedisClient) GetCurrentRate(traceId, originCurrency, destinationCurrency string) (rate float64, expiredAt time.Time) {
	key := strings.Replace(r.svcConfig.RedisConfig.KeyRate, "{origin_currency}", originCurrency, 1)
	key = strings.Replace(key, "{destination_currency}", destinationCurrency, 1)
	fields := []string{"rate", "expired_at"}
	cmd := r.redisClient.HMGet(key, fields...)
	if cmd.Err() != nil {
		log.Err(fmt.Errorf("%v, %v", traceId, cmd.Err()))
		return
	} else if len(fields) == len(cmd.Val()) {
		var err error
		if cmd.Val()[0] != nil {
			rate, err = strconv.ParseFloat(fmt.Sprint(cmd.Val()[0]), 64)
			if err != nil {
				log.Err(fmt.Errorf("%v, %v", traceId, err))
				return
			}
		}
		if cmd.Val()[1] != nil {
			expiredAt, err = time.Parse(time.RFC3339, fmt.Sprint(cmd.Val()[1]))
			if err != nil {
				log.Err(fmt.Errorf("%v, %v", traceId, err))
				return
			}
		}
	}
	return rate, expiredAt
}

func (r RedisClient) SetBookRate(traceId, bookId, originCurrency, destinationCurrency string, rate float64,
	expiredAt time.Time) error {
	key := strings.Replace(r.svcConfig.RedisConfig.KeyBookRate, "{book_id}", bookId, 1)
	hMap := make(map[string]interface{})
	hMap["origin_currency"] = originCurrency
	hMap["destination_currency"] = destinationCurrency
	hMap["rate"] = rate
	hMap["expired_at"] = expiredAt
	boolCmd := r.redisClient.HMSet(key, hMap)
	if boolCmd.Err() != nil {
		log.Err(fmt.Errorf("%v, %v", traceId, boolCmd.Err()))
		return boolCmd.Err()
	}
	boolCmd = r.redisClient.ExpireAt(key, expiredAt)
	if boolCmd.Err() != nil {
		log.Err(fmt.Errorf("%v, %v", traceId, boolCmd.Err()))
		return boolCmd.Err()
	}
	return nil
}
