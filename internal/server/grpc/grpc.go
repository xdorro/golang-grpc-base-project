package grpc

import (
	"sync"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xdorro/golang-micro-base-project/internal/log"
	"github.com/xdorro/golang-micro-base-project/internal/service"
)

type server struct {
	sync.Mutex
	// tlsCredentials     credentials.TransportCredentials
	options            []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
}

// NewGrpcServer returns a Server.
func NewGrpcServer(tlsCredentials credentials.TransportCredentials) Server {
	return &server{
		// tlsCredentials: tlsCredentials,
	}
}

func (s *server) Start(register RegisterFn) *grpc.Server {
	// log gRPC library internals with log
	grpc_zap.ReplaceGrpcLoggerV2(log.Logger())

	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(log.Logger()),
		grpc_recovery.StreamServerInterceptor(),
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(log.Logger()),
		grpc_recovery.UnaryServerInterceptor(),
	}

	options := append(
		s.options,
		WithUnaryServerInterceptors(unaryInterceptors...),
		WithStreamServerInterceptors(streamInterceptors...),
	)

	srv := grpc.NewServer(options...)
	svc := service.NewService()
	s.Lock()
	register(srv, svc)
	s.Unlock()

	return srv
}

func (s *server) Stop() error {
	return nil
}
