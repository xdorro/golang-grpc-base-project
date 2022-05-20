package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"strings"
	"sync"

	"github.com/elastic/gmux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xdorro/golang-grpc-base-project/internal/log"
	grpcS "github.com/xdorro/golang-grpc-base-project/internal/server/grpc"
	httpS "github.com/xdorro/golang-grpc-base-project/internal/server/http"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

type Server interface {
	// Run the server
	Run() error
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

	grpcS := grpcS.NewGrpcServer()
	s.grpc = grpcS.Start(func(srv *grpc.Server, svc service.Service) {
		userpb.RegisterUserServiceServer(srv, svc)
	})

	httpS := httpS.NewHttpServer(s.address, tlsCredentials)
	s.http = httpS.Start(func(srv *runtime.ServeMux, conn *grpc.ClientConn) {
		if err = userpb.RegisterUserServiceHandler(s.ctx, srv, conn); err != nil {
			log.Panicf("proto.RegisterUserServiceHandler(): %w", err)
		}
	})

	// mixed := s.MixedHandler()
	// return http.ListenAndServeTLS(s.address, cert, key, mixed)

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

// MixedHandler returns a handler that runs both http and grpc handlers.
func (s *server) MixedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("content-type"), "application/grpc") {
			s.grpc.ServeHTTP(w, r)
			return
		}

		s.http.ServeHTTP(w, r)
	})
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
