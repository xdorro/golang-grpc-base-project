package redis

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ProviderSet is redis providers.
var ProviderSet = wire.NewSet(NewRedis)

func NewRedis(ctx context.Context, log *zap.Logger) redis.UniversalClient {
	log.Info("Connecting to redis...")

	redisURL := strings.Split(strings.Trim(viper.GetString("REDIS_URL"), " "), ",")
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    redisURL,
		Password: viper.GetString("REDIS_PASSWORD"), // no password set
		DB:       viper.GetInt("REDIS_DB"),          // use default DB

		PoolSize: 1000,
	})

	redisCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := rdb.Ping(redisCtx).Err(); err != nil {
		log.Panic("rdb.Ping()", zap.Error(err))
	}

	log.Info("Redis connected")

	return rdb
}
