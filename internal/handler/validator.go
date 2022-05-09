package handler

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "github.com/xdorro/proto-base-project/protos/v1/auth"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
)

// ValidateError validate payload error if not nil
func (h *Handler) ValidateError(err error) error {
	if err != nil {
		if e, ok := err.(validation.Errors); ok {
			for name, value := range e {
				return status.New(codes.InvalidArgument, fmt.Sprintf("%v: %v", name, value)).Err()
			}
		}

		return status.New(codes.InvalidArgument, err.Error()).Err()
	}

	return nil
}

// ValidateLoginRequest validate login request
func (h *Handler) ValidateLoginRequest(req *authpb.LoginRequest) error {
	err := validation.ValidateStruct(req,
		// Validate phone
		validation.Field(&req.Email,
			validation.Required,
			is.Email,
		),
		// Validate password
		validation.Field(&req.Password,
			validation.Required,
			validation.Length(5, 0),
		),
	)

	return h.ValidateError(err)
}

// ValidateTokenRequest validate token request
func (h *Handler) ValidateTokenRequest(req *authpb.TokenRequest) error {
	err := validation.ValidateStruct(req,
		// Validate token
		validation.Field(&req.Token,
			validation.Required,
			validation.Length(5, 0),
		),
	)

	return h.ValidateError(err)
}

// ValidateCreateUserRequest validate create user request
func (h *Handler) ValidateCreateUserRequest(req *userpb.CreateUserRequest) error {
	err := validation.ValidateStruct(req,
		// Validate name
		validation.Field(&req.Name,
			validation.Required,
			validation.Length(5, 0),
		),
		// Validate email
		validation.Field(&req.Email,
			validation.Required,
			is.Email,
		),
		// Validate password
		validation.Field(&req.Password,
			validation.Required,
			validation.Length(5, 0),
		),
		// Validate roles
		// validation.Field(&req.Role,
		// 	validation.Required.When(req.GetRole() != ""),
		// ),
	)

	return h.ValidateError(err)
}

// ValidateUpdateUserRequest validate update user request
func (h *Handler) ValidateUpdateUserRequest(req *userpb.UpdateUserRequest) error {
	err := validation.ValidateStruct(req,
		// Validate name
		validation.Field(&req.Name,
			validation.Required,
			validation.Length(5, 0),
		),
		// Validate email
		validation.Field(&req.Email,
			validation.Required,
			is.Email,
		),
		// Validate roles
		// validation.Field(&req.Role,
		// 	validation.Required.When(req.GetRole() != ""),
		// ),
	)

	return h.ValidateError(err)
}
