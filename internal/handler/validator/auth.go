package validator_handler

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/xdorro/golang-grpc-base-project/api/proto/auth"
)

func (val *Validator) ValidateLoginRequest(in *auth_proto.LoginRequest) error {
	err := validation.ValidateStruct(in,
		// Validate phone
		validation.Field(&in.Email,
			validation.Required,
			is.Email,
		),
		// Validate password
		validation.Field(&in.Password,
			validation.Required,
			validation.Length(5, 0),
		),
	)

	return val.ValidateError(err)
}

func (val *Validator) ValidateTokenRequest(in *auth_proto.TokenRequest) error {
	err := validation.ValidateStruct(in,
		// Validate token
		validation.Field(&in.Token,
			validation.Required,
			validation.Length(5, 0),
		),
	)

	return val.ValidateError(err)
}
