package main

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/xdorro/golang-grpc-base-project/config"
	"github.com/xdorro/golang-grpc-base-project/internal/log"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.InitConfig()

	host := fmt.Sprintf("localhost:%s", viper.GetString("APP_PORT"))
	log.Infof("Starting https://%s", host)
	srv := server.NewServer(host)

	go func(srv server.Server) {
		if err := srv.Run(); err != nil {
			log.Panicf("error running app: %v", err)
		}
	}(srv)

	// wait for termination signal and register client & http server clean-up operations
	wait := utils.GracefulShutdown(ctx, utils.DefaultShutdownTimeout, map[string]utils.Operation{})
	<-wait
}
