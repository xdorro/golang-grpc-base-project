package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	user_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/user"
)

func (val *Validator) ValidateCreateUserRequest(in *user_proto.CreateUserRequest) error {
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
		),
		// Validate password
		validation.Field(&in.Password,
			validation.Required,
			validation.Length(5, 0),
		),
		// Validate roles
		validation.Field(&in.Role,
			validation.Required.When(in.GetRole() != ""),
		),
	)

	return ValidateError(err)
}

func (val *Validator) ValidateUpdateUserRequest(in *user_proto.UpdateUserRequest) error {
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
		),
		// Validate roles
		validation.Field(&in.Role,
			validation.Required.When(in.GetRole() != ""),
		),
	)

	return ValidateError(err)
}
