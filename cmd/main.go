package main

import (
	"context"
	"runtime"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/config"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/pkg/client"
	"github.com/xdorro/golang-grpc-base-project/pkg/redis"
)

const (
	defaultShutdownTimeout = 10 * time.Second
)

func init() {
	// Load config environment
	config.NewConfig()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// declare new option
	opts := common.NewOption(ctx)
	opts.Log.Info(viper.GetString("APP_NAME"),
		zap.String("app-version", viper.GetString("APP_VERSION")),
		zap.String("go-version", runtime.Version()),
	)

	// declare new client
	client.NewClient(opts)
	// declare new redis
	redis.NewRedis(opts)

	// create new server
	srv, err := server.NewServer(opts)
	if err != nil {
		opts.Log.Fatal("server.NewServer()", zap.Error(err))
	}

	// wait for termination signal and register database & http server clean-up operations
	wait := gracefulShutdown(opts, defaultShutdownTimeout, map[string]operation{
		"server": func(ctx context.Context) error {
			return srv.Close()
		},
		"client": func(ctx context.Context) error {
			return opts.Client.Close()
		},
		"redis": func(ctx context.Context) error {
			return opts.Redis.Close()
		},
		"logger": func(ctx context.Context) error {
			return opts.Log.Sync()
		},
	})

	<-wait
}
