package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/wangn-tech/bookstore-go/internal/app/config"
	"github.com/wangn-tech/bookstore-go/pkg/logger"
)

var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() {
	conf := config.AppConf.Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	RedisClient = redisClient
	logger.Log.Info("Redis connected successfully")
}
