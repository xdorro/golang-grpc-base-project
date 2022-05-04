package redis

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

// ProviderRedisSet is redis providers.
var ProviderRedisSet = wire.NewSet(NewRedis)

// NewRedis is new redis.
func NewRedis(ctx context.Context) redis.UniversalClient {
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
