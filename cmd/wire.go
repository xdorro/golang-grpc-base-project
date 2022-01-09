//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"context"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/internal/handler"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/internal/validator"
	"github.com/xdorro/golang-grpc-base-project/pkg/client"
	"github.com/xdorro/golang-grpc-base-project/pkg/redis"
)

func initializeServer(ctx context.Context, log *zap.Logger) (*server.Server, error) {
	wire.Build(
		client.ProviderSet,
		repo.ProviderSet,
		redis.ProviderSet,
		handler.ProviderSet,
		validator.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
	)

	return &server.Server{}, nil
}
