package main

import (
	"context"
	"runtime"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/internal/config"
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

	// wait for termination signal and register database & http server clean-up operations
	wait := gracefulShutdown(opts, defaultShutdownTimeout, map[string]operation{
		"logger": func(ctx context.Context) error {
			return opts.Log.Sync()
		},
	})

	<-wait
}
