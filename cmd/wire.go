//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"

	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

func initServer() server.IServer {
	wire.Build(
		service.ProviderServiceSet,
		server.ProviderServerSet,
	)

	return &server.Server{}
}
