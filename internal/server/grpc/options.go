package grpc

import (
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

// RegisterFn defines the method to register a server.
type RegisterFn func(*grpc.Server, service.IService)

// IServer interface represents a rpc server.
type IServer interface {
	AddOptions(options ...grpc.ServerOption)
	AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)
	AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor)
	Server() *grpc.Server
	Close()
}

func (s *server) AddOptions(options ...grpc.ServerOption) {
	s.Lock()
	s.options = append(s.options, options...)
	s.Unlock()
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
