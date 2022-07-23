//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"

	"github.com/xdorro/golang-grpc-base-project/internal/server"
)

func initServer(opts ...server.Option) server.IServer {
	wire.Build(
		server.ProviderServerSet,
	)

	return &server.Server{}
}
