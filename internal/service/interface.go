package service

import (
	"google.golang.org/grpc"
)

// IService is the interface for the service
type IService interface {
	RegisterServiceServer(grpcServer *grpc.Server)
}
