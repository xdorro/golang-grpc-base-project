package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/elastic/gmux"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	grpcS "github.com/xdorro/golang-grpc-base-project/internal/server/grpc"
	httpS "github.com/xdorro/golang-grpc-base-project/internal/server/http"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

type Server struct {
	sync.Mutex
	ctx  context.Context
	grpc *grpc.Server
}

func NewServer(ctx context.Context, service service.IService) IServer {
	s := &Server{
		ctx:  ctx,
		grpc: grpcS.NewGrpcServer(service, grpcS.RegisterGRPC),
	}

	return s
}

func (s *Server) Run() error {
	host := fmt.Sprintf(":%s", viper.GetString("APP_PORT"))
	log.Infof("Starting http://localhost%s", host)

	srv := &http.Server{
		Addr:    host,
		Handler: httpS.NewHttpServer(host, httpS.RegisterHTTP),
	}

	// Configure the IServer with gmux. The returned net.Listener will receive gRPC connections,
	// while all other requests will be handled by s.Handler.
	grpcListener, err := gmux.ConfigureServer(srv, nil)
	if err != nil {
		log.Fatalf("error configuring server: %v", err)
	}

	go func() {
		if err = s.grpc.Serve(grpcListener); err != nil {
			log.Fatalf("grpc serve error: %v", err)
		}
	}()

	return srv.ListenAndServe()
}

func (s *Server) Close() error {
	s.Lock()
	defer s.Unlock()
	s.grpc.GracefulStop()

	return nil
}
