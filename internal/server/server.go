package server

import (
	"context"
	"fmt"

	"github.com/google/wire"

	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

// ProviderServerSet is Server providers.
var ProviderServerSet = wire.NewSet(NewServer)

// IServer is the interface that must be implemented by a server.
type IServer interface {
	Run() error
	Close() error
}

// Server is a server struct.
type Server struct {
	ctx     context.Context
	name    string
	version string
	port    int
}

// NewServer creates a new server.
func NewServer(opts ...Option) IServer {
	s := Server{}

	// Loop through each option
	for _, opt := range opts {
		opt(&s)
	}

	return s
}

// Run runs the server.
func (s Server) Run() error {
	log.Info().
		Str("app-name", s.name).
		Str("app-version", s.version).
		Int("app-port", s.port).
		Msg("Config loaded")

	host := fmt.Sprintf(":%d", s.port)
	log.Infof("Starting application http://localhost%s", host)

	return nil
}

// Close closes the server.
func (s Server) Close() error {
	return nil
}
