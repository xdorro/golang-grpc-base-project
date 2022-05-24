package grpc

import (
	authpb "github.com/xdorro/proto-base-project/protos/v1/auth"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

// RegisterFn defines the method to register a server.
type RegisterFn func(*grpc.Server, service.IService)

// IServer interface represents a rpc server.
type IServer interface {
	AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)
	AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor)
	Server() *grpc.Server
	Close()
}

func (s *server) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	s.Lock()
	s.streamInterceptors = append(s.streamInterceptors, interceptors...)
	s.Unlock()
}

func (s *server) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	s.Lock()
	s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)
	s.Unlock()
}

func RegisterGRPC(srv *grpc.Server, svc service.IService) {
	userpb.RegisterUserServiceServer(srv, svc)
	authpb.RegisterAuthServiceServer(srv, svc)
}
