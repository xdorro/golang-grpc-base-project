package handler

import (
	authpb "github.com/xdorro/base-project-proto/protos/v1/auth"
	userpb "github.com/xdorro/base-project-proto/protos/v1/user"
)

// IHandler is the interface for the server
type IHandler interface {
	IValidateHandler
}

// IValidateHandler is the interface for the validation handler
type IValidateHandler interface {
	ValidateError(err error) error

	ValidateLoginRequest(req *authpb.LoginRequest) error
	ValidateTokenRequest(req *authpb.TokenRequest) error

	ValidateCreateUserRequest(req *userpb.CreateUserRequest) error
	ValidateUpdateUserRequest(req *userpb.UpdateUserRequest) error
}
