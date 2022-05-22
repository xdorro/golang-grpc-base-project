package grpc

import (
	authpb "github.com/xdorro/proto-base-project/protos/v1/auth"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

func RegisterGRPC(srv *grpc.Server, svc service.Service) {
	userpb.RegisterUserServiceServer(srv, svc)
	authpb.RegisterAuthServiceServer(srv, svc)
}
