package service

import (
	authpb "github.com/xdorro/proto-base-project/protos/v1/auth"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
)

type Service interface {
	userpb.UserServiceServer
	authpb.AuthServiceServer
}

// service is service struct.
type service struct {
	userpb.UnimplementedUserServiceServer
	authpb.UnimplementedAuthServiceServer
}

// NewService creates a new service.
func NewService() Service {
	return &service{}
}
