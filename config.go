package main

import (
	"flag"
	"fmt"
	"go-forex/model"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
)

func GetConfig() (svcConfig model.SvcConfig) {
	isLocalDev := flag.Bool("local", false, "=(true/false)")
	flag.Parse()
	if *isLocalDev {
		if err := godotenv.Load(".env"); err != nil {
			panic(err)
		}
	}
	return getEnv()
}

func getEnv() model.SvcConfig {
	svcConfig := model.SvcConfig{
		URLGetRate:       os.Getenv("URL_GET_RATES"),
		URLGetRateMethod: os.Getenv("URL_GET_RATES_METHOD"),
		RedisConfig: model.RedisConfig{
			Host:        os.Getenv("REDIS_HOST"),
			KeyRate:     os.Getenv("REDIS_KEY_RATE"),
			KeyBookRate: os.Getenv("REDIS_KEY_BOOK_RATE"),
		},
		DatabaseConfig: model.DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
		},
	}
	err := validator.New().Struct(&svcConfig)
	if err != nil {
		panic(err)
	}
	return svcConfig
}

func redisConnect(svcConfig model.SvcConfig) *redis.Client {
	fmt.Printf("Connecting to Redis : '%v'\n", svcConfig.RedisConfig.Host)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     svcConfig.RedisConfig.Host,
		Password: "", // no password set
		DB:       0,  // use default DB)
	})
	result := redisClient.Ping()
	if result.Err() != nil {
		fmt.Println(result.Err(), "Error connnect to redis : "+result.Err().Error())
		time.Sleep(7 * time.Second)
		os.Exit(1)
	}
	fmt.Printf("Connected to Redis : '%v'\n", svcConfig.RedisConfig.Host)
	return redisClient
}

func databaseConnect(svcConfig model.SvcConfig) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     svcConfig.DatabaseConfig.Host,
		User:     svcConfig.DatabaseConfig.User,
		Password: svcConfig.DatabaseConfig.Password,
		Database: svcConfig.DatabaseConfig.Database,
	})
	var check int
	if _, err := db.QueryOne(pg.Scan(&check), "SELECT 1"); err != nil {
		panic(err)
	}
	return db
}
