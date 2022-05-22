package server

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func (s *server) AddGRPCServer(grpc *grpc.Server) {
	s.mu.Lock()
	s.grpc = grpc
	s.mu.Unlock()
}

func (s *server) AddHTTPServer(http *runtime.ServeMux) {
	s.mu.Lock()
	s.http = http
	s.mu.Unlock()
}
