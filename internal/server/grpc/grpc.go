package grpc

import (
	"sync"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

type server struct {
	sync.Mutex
	options            []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
}

// NewGrpcServer returns a Server.
func NewGrpcServer() Server {
	return &server{}
}

func (s *server) Start(register RegisterFn) *grpc.Server {
	// log gRPC library internals with log
	// grpc_zap.ReplaceGrpcLoggerV2(log.Logger())

	s.AddStreamInterceptors(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		// grpc_zap.StreamServerInterceptor(log.Logger()),
		grpc_recovery.StreamServerInterceptor(),
	)

	s.AddUnaryInterceptors(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		// grpc_zap.UnaryServerInterceptor(log.Logger()),
		grpc_recovery.UnaryServerInterceptor(),
	)

	s.AddOptions(
		WithUnaryServerInterceptors(s.unaryInterceptors...),
		WithStreamServerInterceptors(s.streamInterceptors...),
	)

	srv := grpc.NewServer(s.options...)
	svc := service.NewService()
	// s.Lock()
	// register(srv, svc)
	// s.Unlock()
	userpb.RegisterUserServiceServer(srv, svc)

	return srv
}

func (s *server) Stop() error {
	return nil
}
