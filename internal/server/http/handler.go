package http

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/xdorro/proto-base-project/protos/v1/auth"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/log"
)

func RegisterHTTP(srv *runtime.ServeMux, conn *grpc.ClientConn) {
	ctx := context.Background()

	if err := userpb.RegisterUserServiceHandler(ctx, srv, conn); err != nil {
		log.Panicf("proto.RegisterUserServiceHandler(): %w", err)
	}

	if err := authpb.RegisterAuthServiceHandler(ctx, srv, conn); err != nil {
		log.Panicf("proto.RegisterUserServiceHandler(): %w", err)
	}
}
