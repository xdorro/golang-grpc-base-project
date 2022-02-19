package validator_handler

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/proto/permission"
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

	return val.ValidateError(err)
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

	return val.ValidateError(err)
}

func (val *Validator) ValidateListPermissions(list []string) ([]*ent.Permission, error) {
	permissions := make([]*ent.Permission, 0)

	if len(list) > 0 {
		permissions = val.repo.FindAllPermissionBySlugs(list)

		for _, slug := range list {
			if err := val.hasSlugInPermissions(permissions, slug); err != nil {
				return nil, err
			}
		}
	}

	return permissions, nil
}

func (val *Validator) hasSlugInPermissions(permissions []*ent.Permission, slug string) error {
	for _, permission := range permissions {
		if strings.EqualFold(slug, permission.Slug) {
			return nil
		}
	}

	return status.New(codes.InvalidArgument, fmt.Sprintf("permission: %s doesn't exist", slug)).Err()
}
