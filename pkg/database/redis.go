package database

import (
	"context"
	"fmt"
	"log"

	"easyfind/internal/config"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() {
	cfg := config.AppConfigData.Database.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}
}
