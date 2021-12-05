package authservice

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/kucow/golang-grpc-base/internal/common"
	authproto "github.com/kucow/golang-grpc-base/pkg/proto/v1/auth"
)

func (svc *AuthService) validateLoginRequest(in *authproto.LoginRequest) error {
	err := validation.ValidateStruct(in,
		// Validate email
		validation.Field(&in.Email,
			validation.Required,
			is.Email,
			validation.Length(5, 0),
		),
		// Validate password
		validation.Field(&in.Password,
			validation.Required,
			validation.Length(5, 0),
		),
	)

	return common.ValidateError(err)
}

func (svc *AuthService) validateTokenRequest(in *authproto.TokenRequest) error {
	err := validation.ValidateStruct(in,
		// Validate token
		validation.Field(&in.Token,
			validation.Required,
			validation.Length(5, 0),
		),
	)

	return common.ValidateError(err)
}
