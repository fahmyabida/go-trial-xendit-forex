package main

import (
	http_client "go-forex/http/client"
	http_server "go-forex/http/server"
	"go-forex/redis_jobs"
	"go-forex/repo"
	"go-forex/usecase"
)

func RunApplication() {
	svcConfig := GetConfig()

	dbConn := databaseConnect(svcConfig)
	redisConn := redisConnect(svcConfig)

	logBookingRateRepo := repo.NewLogBookingRateRepo(dbConn)
	redisClient := redis_jobs.NewRedisClient(redisConn, svcConfig)
	httpClient := http_client.NewHttpClient(svcConfig)

	getRateUsecase := usecase.NewGetRateUsecase(httpClient, redisClient, logBookingRateRepo)
	serverHttp := http_server.NewServerHttp(getRateUsecase)

	serverHttp.Init()
}
