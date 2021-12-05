package main

import (
	"context"
	"runtime"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/internal/config"
	"github.com/kucow/golang-grpc-base/internal/server"
	"github.com/kucow/golang-grpc-base/pkg/client"
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

	// declare client
	client.NewClient(opts)
	// // declare redis
	// redis.NewRedis(opts)

	// start server
	srv, err := server.NewServer(opts)
	if err != nil {
		opts.Log.Fatal("server.NewServer()", zap.Error(err))
	}

	// wait for termination signal and register database & http server clean-up operations
	wait := gracefulShutdown(opts, defaultShutdownTimeout, map[string]operation{
		"client": func(ctx context.Context) error {
			return opts.Client.Close()
		},
		"server": func(ctx context.Context) error {
			return srv.Close()
		},
		"logger": func(ctx context.Context) error {
			return opts.Log.Sync()
		},
	})

	<-wait
}
