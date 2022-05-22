package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"

	"github.com/elastic/gmux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xdorro/golang-grpc-base-project/internal/log"
	grpcS "github.com/xdorro/golang-grpc-base-project/internal/server/grpc"
	httpS "github.com/xdorro/golang-grpc-base-project/internal/server/http"
)

type Server interface {
	// Run the server
	Run() error
	Close() error
}

type server struct {
	mu      sync.Mutex
	address string
	grpc    *grpc.Server
	http    *runtime.ServeMux
	ctx     context.Context
}

func NewServer(address string) Server {
	s := &server{
		ctx:     context.Background(),
		address: address,
	}

	return s
}

func (s *server) Run() error {
	cert := viper.GetString("APP_CERT")
	key := viper.GetString("APP_KEY")
	tlsCredentials, err := LoadTLSCredentials(cert, key)
	if err != nil {
		log.Panicf("cannot load TLS credentials: %s", err)
	}

	newGrpc := grpcS.NewGrpcServer()
	s.grpc = newGrpc.Start(grpcS.RegisterGRPC)

	newHttp := httpS.NewHttpServer(s.address, tlsCredentials)
	s.http = newHttp.Start(httpS.RegisterHTTP)

	srv := &http.Server{
		Addr:    s.address,
		Handler: s.http,
	}

	// Configure the server with gmux. The returned net.Listener will receive gRPC connections,
	// while all other requests will be handled by s.Handler.
	grpcListener, err := gmux.ConfigureServer(srv, nil)
	if err != nil {
		log.Fatalf("error configuring server: %", err)
	}

	go func() {
		if err = s.grpc.Serve(grpcListener); err != nil {
			log.Fatalf("grpc server error: %v", err)
		}
	}()

	return srv.ListenAndServeTLS(cert, key)
}

// LoadTLSCredentials loads TLS credentials from the configuration
func LoadTLSCredentials(cert, key string) (
	credentials.TransportCredentials, error,
) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates:       []tls.Certificate{serverCert},
		ClientAuth:         tls.RequireAndVerifyClientCert,
		InsecureSkipVerify: true,
	}

	return credentials.NewTLS(config), nil
}

func (s *server) Close() error {
	s.grpc.GracefulStop()

	return nil
}
