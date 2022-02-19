package validator_handler

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/proto/role"
)

func (val *Validator) ValidateCreateRoleRequest(in *role_proto.CreateRoleRequest) error {
	err := validation.ValidateStruct(in,
		// Validate name
		validation.Field(&in.Name,
			validation.Required,
			validation.Length(3, 0),
		),
		// Validate slug
		validation.Field(&in.Slug,
			validation.Required,
			validation.Length(3, 0),
		),
		// Validate permissions
		validation.Field(&in.Permissions,
			validation.Required.When(in.GetPermissions() != nil),
		),
	)

	return val.ValidateError(err)
}

func (val *Validator) ValidateUpdateRoleRequest(in *role_proto.UpdateRoleRequest) error {
	err := validation.ValidateStruct(in,
		// Validate id
		validation.Field(&in.Id,
			validation.Required,
			is.Int,
		),
		// Validate name
		validation.Field(&in.Name,
			validation.Required,
			validation.Length(3, 0),
		),
		// Validate slug
		validation.Field(&in.Slug,
			validation.Required,
			validation.Length(3, 0),
		),
		// Validate permissions
		validation.Field(&in.Permissions,
			validation.Required.When(in.GetPermissions() != nil),
		),
	)

	return val.ValidateError(err)
}

func (val *Validator) ValidateRole(slug string) (*ent.Role, error) {
	role, err := val.repo.FindRoleBySlug(slug)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, fmt.Sprintf("role: %s doesn't exist", slug)).Err()
	}

	return role, nil
}
