package service

import (
	authpb "github.com/xdorro/proto-base-project/protos/v1/auth"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
	"google.golang.org/grpc"
)

// IService is the interface for the service
type IService interface {
	RegisterServiceServer(grpcServer *grpc.Server)

	authpb.AuthServiceServer
	userpb.UserServiceServer
}
