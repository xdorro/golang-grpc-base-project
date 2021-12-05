package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	authproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/auth"
)

func (val *Validator) ValidateLoginRequest(in *authproto.LoginRequest) error {
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

	return ValidateError(err)
}

func (val *Validator) ValidateTokenRequest(in *authproto.TokenRequest) error {
	err := validation.ValidateStruct(in,
		// Validate token
		validation.Field(&in.Token,
			validation.Required,
			validation.Length(5, 0),
		),
	)

	return ValidateError(err)
}
