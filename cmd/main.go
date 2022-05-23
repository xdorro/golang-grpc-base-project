package main

import (
	"context"

	"github.com/xdorro/golang-grpc-base-project/config"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/pkg/log"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init config
	config.InitConfig()

	// New server
	srv := initServer(ctx)

	// Run server
	go func(srv server.IServer) {
		if err := srv.Run(); err != nil {
			log.Panicf("error running app: %v", err)
		}
	}(srv)

	// wait for termination signal and register client & http server clean-up operations
	wait := utils.GracefulShutdown(ctx, utils.DefaultShutdownTimeout, map[string]utils.Operation{
		"server": func(ctx context.Context) error {
			return srv.Close()
		},
	})
	<-wait
}
