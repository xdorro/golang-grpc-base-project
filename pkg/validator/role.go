package validator

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	roleproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/role"
)

func (val *Validator) ValidateCreateRoleRequest(in *roleproto.CreateRoleRequest) error {
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
			validation.Each(IsULID),
		),
	)

	return ValidateError(err)
}

func (val *Validator) ValidateUpdateRoleRequest(in *roleproto.UpdateRoleRequest) error {
	err := validation.ValidateStruct(in,
		// Validate id
		validation.Field(&in.Id,
			validation.Required,
			validation.Length(5, 100),
			IsULID,
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
			validation.Each(IsULID),
		),
	)

	return ValidateError(err)
}

func (val *Validator) ValidateListRoles(list []uint64) ([]*ent.Role, error) {
	roles := make([]*ent.Role, 0)

	if len(list) > 0 {
		for _, id := range list {
			r, err := val.persist.FindRoleByID(id)
			if err != nil {
				return nil, status.New(codes.InvalidArgument, fmt.Sprintf("role: %s doesn't exist", id)).Err()
			}

			roles = append(roles, r)
		}
	}

	return roles, nil
}
