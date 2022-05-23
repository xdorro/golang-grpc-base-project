package service

import (
	"github.com/google/wire"
	authpb "github.com/xdorro/proto-base-project/protos/v1/auth"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
)

// ProviderServiceSet is Server providers.
var ProviderServiceSet = wire.NewSet(NewService)

type IService interface {
	userpb.UserServiceServer
	authpb.AuthServiceServer
}
