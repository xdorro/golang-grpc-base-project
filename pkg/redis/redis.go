package redis

import (
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

var rdb redis.UniversalClient

func NewRedis() redis.UniversalClient {
	logger.Info("Connecting to redis...")

	redisURL := strings.Split(strings.Trim(viper.GetString("REDIS_URL"), " "), ",")
	rdb = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    redisURL,
		Password: viper.GetString("REDIS_PASSWORD"), // no password set
		DB:       viper.GetInt("REDIS_DB"),          // use default DB

		PoolSize:     1000,
		PoolTimeout:  2 * time.Minute,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	})

	if err := rdb.Ping(rdb.Context()).Err(); err != nil {
		logger.Fatal("rdb.Ping()", zap.Error(err))
	}

	logger.Info("redis connected")

	return rdb
}
