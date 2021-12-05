package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	userproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/user"
)

func (val *Validator) ValidateCreateUserRequest(in *userproto.CreateUserRequest) error {
	err := validation.ValidateStruct(in,
		// Validate name
		validation.Field(&in.Name,
			validation.Required,
			validation.Length(5, 0),
		),
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
		// Validate roles
		validation.Field(&in.Roles,
			validation.Required,
			validation.Each(is.Int),
		),
	)

	return ValidateError(err)
}

func (val *Validator) ValidateUpdateUserRequest(in *userproto.UpdateUserRequest) error {
	err := validation.ValidateStruct(in,
		// Validate id
		validation.Field(&in.Id,
			validation.Required,
			is.Int,
		),
		// Validate name
		validation.Field(&in.Name,
			validation.Required,
			validation.Length(5, 0),
		),
		// Validate email
		validation.Field(&in.Email,
			validation.Required,
			is.Email,
			validation.Length(5, 0),
		),
		// Validate roles
		validation.Field(&in.Roles,
			validation.Required,
			validation.Each(is.Int),
		),
	)

	return ValidateError(err)
}
