package permissionservice

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/kucow/golang-grpc-base/internal/common"
	permissionproto "github.com/kucow/golang-grpc-base/pkg/proto/v1/permission"
)

func (svc *PermissionService) validateCreatePermissionRequest(in *permissionproto.CreatePermissionRequest) error {
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

	return common.ValidateError(err)
}

func (svc *PermissionService) validateUpdatePermissionRequest(in *permissionproto.UpdatePermissionRequest) error {
	err := validation.ValidateStruct(in,
		// Validate id
		validation.Field(&in.Id,
			validation.Required,
			validation.Length(5, 100),
			common.IsULID,
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

	return common.ValidateError(err)
}
