package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"

	"github.com/elastic/gmux"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials"

	grpcS "github.com/xdorro/golang-grpc-base-project/internal/server/grpc"
	httpS "github.com/xdorro/golang-grpc-base-project/internal/server/http"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

type Server struct {
	sync.Mutex
	ctx  context.Context
	grpc grpcS.IServer
}

func NewServer(ctx context.Context, service service.IService) IServer {
	s := &Server{
		ctx:  ctx,
		grpc: grpcS.NewGrpcServer(service, grpcS.RegisterGRPC),
	}

	return s
}

func (s *Server) Run() error {
	cert := viper.GetString("APP_CERT")
	key := viper.GetString("APP_KEY")
	tlsCredentials, err := LoadTLSCredentials(cert, key)
	if err != nil {
		log.Panicf("cannot load TLS credentials: %s", err)
	}

	host := fmt.Sprintf("localhost:%s", viper.GetString("APP_PORT"))
	log.Infof("Starting https://%s", host)

	newHttp := httpS.NewHttpServer(host, tlsCredentials)
	srv := &http.Server{
		Addr:    host,
		Handler: newHttp.Start(httpS.RegisterHTTP),
	}

	// Configure the IServer with gmux. The returned net.Listener will receive gRPC connections,
	// while all other requests will be handled by s.Handler.
	grpcListener, err := gmux.ConfigureServer(srv, nil)
	if err != nil {
		log.Fatalf("error configuring IServer: %v", err)
	}

	go func() {
		if err = s.grpc.Server().Serve(grpcListener); err != nil {
			log.Fatalf("grpc IServer error: %v", err)
		}
	}()

	return srv.ListenAndServeTLS(cert, key)
}

// LoadTLSCredentials loads TLS credentials from the configuration
func LoadTLSCredentials(cert, key string) (
	credentials.TransportCredentials, error,
) {
	// Load IServer's certificate and private key
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

func (s *Server) Close() error {
	s.grpc.Close()

	return nil
}
