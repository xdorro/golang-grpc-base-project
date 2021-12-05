package redis

import (
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base/internal/common"
)

func NewRedis(opts *common.Option) {
	opts.Log.Info("Connecting to Redis...")

	redisURL := strings.Split(strings.Trim(viper.GetString("REDIS_URL"), " "), ",")
	opts.Redis = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    redisURL,
		Password: viper.GetString("REDIS_PASSWORD"), // no password set
		DB:       viper.GetInt("REDIS_DB"),          // use default DB

		PoolSize:     1000,
		PoolTimeout:  2 * time.Minute,
		IdleTimeout:  10 * time.Minute,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 1 * time.Minute,
	})

	if err := opts.Redis.Ping(opts.Ctx).Err(); err != nil {
		opts.Log.Fatal("rdb.Ping()", zap.Error(err))
	}

	opts.Log.Info("Redis connected")
}
