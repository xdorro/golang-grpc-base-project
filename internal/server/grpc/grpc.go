package grpc

import (
	"context"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tags"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/log"
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
	logger := zerolog.InterceptorLogger(log.Logger)

	serverPayloadDecider := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		return true
	}

	s.AddStreamInterceptors(
		tags.StreamServerInterceptor(tags.WithFieldExtractor(tags.CodeGenRequestFieldExtractor)),
		logging.StreamServerInterceptor(logger),
		logging.PayloadStreamServerInterceptor(logger, serverPayloadDecider),
		recovery.StreamServerInterceptor(),
	)

	s.AddUnaryInterceptors(
		tags.UnaryServerInterceptor(tags.WithFieldExtractor(tags.CodeGenRequestFieldExtractor)),
		logging.UnaryServerInterceptor(logger),
		logging.PayloadUnaryServerInterceptor(logger, serverPayloadDecider),
		recovery.UnaryServerInterceptor(),
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
