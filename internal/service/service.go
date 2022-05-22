package service

import (
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
)

type Service interface {
	userpb.UserServiceServer
}

// service is service struct.
type service struct {
	userpb.UnimplementedUserServiceServer
}

// NewService creates a new service.
func NewService() Service {
	return &service{}
}
