package core

import (
	"context"

	"github.com/xdorro/golang-grpc-base-project/internal/server/grpc"
	"github.com/xdorro/golang-grpc-base-project/internal/server/http"
)

type Option func(*Options)

// Options for app
type Options struct {
	Name    string
	Version string
	Address string

	grpcServer grpc.Server
	httpServer http.Server

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Context: context.Background(),
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service and for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// Name of the service
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// Version of the service
func Version(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

// Address of the service
func Address(address string) Option {
	return func(o *Options) {
		o.Address = address
	}
}
