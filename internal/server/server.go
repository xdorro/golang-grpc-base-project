package server

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/elastic/gmux"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	grpcS "github.com/xdorro/golang-grpc-base-project/internal/server/grpc"
	httpS "github.com/xdorro/golang-grpc-base-project/internal/server/http"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

type Server struct {
	grpc *grpc.Server
	http *http.Server
}

// NewServer creates a new server
func NewServer(service service.IService) IServer {
	host := fmt.Sprintf(":%s", viper.GetString("APP_PORT"))
	log.Infof("Starting application http://localhost%s", host)

	s := &Server{
		grpc: grpcS.NewGrpcServer(service, grpcS.RegisterGRPC),
		http: &http.Server{
			Addr:    host,
			Handler: httpS.NewHttpServer(host, httpS.RegisterHTTP),
		},
	}

	// Configure the IServer with gmux. The returned net.Listener will receive gRPC connections,
	// while all other requests will be handled by s.Handler.
	grpcListener, err := gmux.ConfigureServer(s.http, nil)
	if err != nil {
		log.Fatalf("error configuring server: %v", err)
	}

	// Serve the gRPC server on the gRPC listener.
	go func() {
		if err = s.grpc.Serve(grpcListener); err != nil {
			log.Fatalf("grpc serve error: %v", err)
		}
	}()

	// Serve the http server on the http listener.
	go func() {
		if err = s.http.ListenAndServe(); err != nil {
			log.Fatalf("http serve error: %v", err)
		}
	}()

	// we need a webserver to get the pprof webserver
	if viper.GetBool("APP_DEBUG") {
		go func() {
			log.Infof("Starting pprof http://localhost:6060")

			if err = http.ListenAndServe("localhost:6060", nil); err != nil {
				log.Fatalf("http serve error: %v", err)
			}
		}()
	}

	return s
}

// Close closes the server
func (s *Server) Close() error {
	s.grpc.GracefulStop()

	return nil
}
