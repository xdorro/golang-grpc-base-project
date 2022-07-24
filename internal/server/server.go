package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/wire"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
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

	http *http.Server
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

type UserService struct {
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

	// create new mux server
	mux := http.NewServeMux()

	// create new service handler
	service.NewService(mux)

	s.http = &http.Server{
		Addr: host,
		// Use h2c so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	return s.http.ListenAndServe()
}

// Close closes the server.
func (s Server) Close() error {
	return s.http.Close()
}
