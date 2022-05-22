package http

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// RegisterFn defines the method to register a server@@.
type RegisterFn func(*runtime.ServeMux, *grpc.ClientConn)

// Server interface represents a rpc server@@.
type Server interface {
	Start(register RegisterFn) *runtime.ServeMux
}
