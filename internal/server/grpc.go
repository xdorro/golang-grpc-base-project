package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

// newGRPCServer creates a new grpc server
func (s *Server) newGRPCServer(tlsCredentials credentials.TransportCredentials, service service.Service) {
	s.mu.Lock()
	defer s.mu.Unlock()

	streamChain := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_recovery.StreamServerInterceptor(),
	}

	unaryChain := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_recovery.UnaryServerInterceptor(),
	}

	// log payload if enabled
	// if viper.GetBool("LOG_PAYLOAD") {
	// 	alwaysLoggingDeciderServer := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
	// 		return true
	// 	}
	//
	// 	streamChain = append(streamChain, grpc_zap.PayloadStreamServerInterceptor(log.NewLogger(), alwaysLoggingDeciderServer))
	// 	unaryChain = append(unaryChain, grpc_zap.PayloadUnaryServerInterceptor(log.NewLogger(), alwaysLoggingDeciderServer))
	// }

	// register grpc service Server
	s.grpcServer = grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc_middleware.WithStreamServerChain(streamChain...),
		grpc_middleware.WithUnaryServerChain(unaryChain...),
	)

	// register service to grpc server
	userpb.RegisterUserServiceServer(s.grpcServer, service)
}
