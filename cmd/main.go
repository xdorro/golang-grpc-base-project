package main

import (
	"context"

	"github.com/xdorro/golang-grpc-base-project/config"
	"github.com/xdorro/golang-grpc-base-project/internal/log"
	"github.com/xdorro/golang-grpc-base-project/internal/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.InitConfig()

	log.Infof("hello world")

	// wait for termination signal and register client & http server clean-up operations
	wait := utils.GracefulShutdown(ctx, utils.DefaultShutdownTimeout, map[string]utils.Operation{})

	<-wait
}
