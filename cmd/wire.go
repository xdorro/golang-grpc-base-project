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
)

func initializeServer(ctx context.Context, log *zap.Logger) server.IServer {
	wire.Build(
		repo.ProviderRepoSet,
		handler.ProviderHandlerSet,
		service.ProviderServiceSet,
		server.ProviderServerSet,
	)

	return &server.Server{}
}
