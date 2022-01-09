package validator

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	permission_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/permission"
)

func (val *Validator) ValidateCreatePermissionRequest(in *permission_proto.CreatePermissionRequest) error {
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
	)

	return ValidateError(err)
}

func (val *Validator) ValidateUpdatePermissionRequest(in *permission_proto.UpdatePermissionRequest) error {
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
	)

	return ValidateError(err)
}

func (val *Validator) ValidateListPermissions(list []string) ([]*ent.Permission, error) {
	permissions := make([]*ent.Permission, 0)

	if len(list) > 0 {
		for _, slug := range list {
			p, err := val.repo.FindPermissionBySlug(slug)
			if err != nil {
				return nil, status.New(codes.InvalidArgument, fmt.Sprintf("permission: %s doesn't exist", slug)).Err()
			}

			permissions = append(permissions, p)
		}
	}

	return permissions, nil
}
