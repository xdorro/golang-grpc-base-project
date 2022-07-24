package main

import (
	"context"
	"net/http"

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
	srv := initServer(
		server.WithContext(ctx),
	)

	// Run server
	go func(srv server.IServer) {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server.Run() error : %v", err)
			return
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
