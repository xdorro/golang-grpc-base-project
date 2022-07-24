package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/google/wire"
	"github.com/spf13/viper"
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
	mu      sync.Mutex
	ctx     context.Context
	name    string
	version string
	port    int
	debug   bool

	http *http.Server
}

// NewServer creates a new server.
func NewServer(opts ...Option) IServer {
	s := &Server{
		ctx:     context.Background(),
		name:    viper.GetString("APP_NAME"),
		version: viper.GetString("APP_VERSION"),
		port:    viper.GetInt("APP_PORT"),
		debug:   viper.GetBool("APP_DEBUG"),
	}

	// Loop through each option
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Run runs the server.
func (s *Server) Run() error {
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
	s.mu.Lock()
	service.NewService(mux)
	s.mu.Unlock()

	// we need a webserver to get the pprof webserver
	if s.debug {
		go func() {
			log.Infof("Starting pprof http://localhost:6060")

			if err := http.ListenAndServe("localhost:6060", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("http serve error: %v", err)
			}
		}()
	}

	// create new http server
	s.http = &http.Server{
		Addr: host,
		// Use h2c, so we can serve HTTP/2 without TLS.
		Handler: h2c.NewHandler(
			newCORS().Handler(mux),
			&http2.Server{},
		),
		ReadHeaderTimeout: time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      1 * time.Minute,
		MaxHeaderBytes:    8 * 1024, // 8KiB
	}

	// Serve the http server on the http listener.
	return s.http.ListenAndServe()
}

// Close closes the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	return s.http.Shutdown(ctx)
}
