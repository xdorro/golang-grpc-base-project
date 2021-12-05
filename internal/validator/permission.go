package validator

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/spf13/cast"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	permissionproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/permission"
)

func (val *Validator) ValidateCreatePermissionRequest(in *permissionproto.CreatePermissionRequest) error {
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

func (val *Validator) ValidateUpdatePermissionRequest(in *permissionproto.UpdatePermissionRequest) error {
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
		for _, id := range list {
			p, err := val.persist.FindPermissionByID(cast.ToUint64(id))
			if err != nil {
				return nil, status.New(codes.InvalidArgument, fmt.Sprintf("Permission: %s doesn't exist", id)).Err()
			}

			permissions = append(permissions, p)
		}
	}

	return permissions, nil
}
