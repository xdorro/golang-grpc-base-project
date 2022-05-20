package grpc

import (
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

// RegisterFn defines the method to register a server.
type RegisterFn func(*grpc.Server, service.Service)

// Server interface represents a rpc server.
type Server interface {
	AddOptions(options ...grpc.ServerOption)
	AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)
	AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor)
	Start(register RegisterFn) *grpc.Server
}

func (s *server) AddOptions(options ...grpc.ServerOption) {
	s.options = append(s.options, options...)
}

func (s *server) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	s.streamInterceptors = append(s.streamInterceptors, interceptors...)
}

func (s *server) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)
}
