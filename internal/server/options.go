package server

import (
	"github.com/google/wire"
)

// ProviderServerSet is Server providers.
var ProviderServerSet = wire.NewSet(NewServer)

type IServer interface {
	// Run the Server
	Run() error
	Close() error
}
