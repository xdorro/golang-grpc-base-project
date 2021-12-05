package roleservice

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/pkg/ent"
	roleproto "github.com/kucow/golang-grpc-base/pkg/proto/v1/role"
)

func (svc *RoleService) validateCreateRoleRequest(in *roleproto.CreateRoleRequest) error {
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
			validation.Each(common.IsULID),
		),
	)

	return common.ValidateError(err)
}

func (svc *RoleService) validateUpdateRoleRequest(in *roleproto.UpdateRoleRequest) error {
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
		// Validate permissions
		validation.Field(&in.Permissions,
			validation.Required.When(in.GetPermissions() != nil),
			validation.Each(common.IsULID),
		),
	)

	return common.ValidateError(err)
}

func (svc *RoleService) validateListPermissions(list []string) ([]*ent.Permission, error) {
	permissions := make([]*ent.Permission, 0)

	if len(list) > 0 {
		for _, id := range list {
			p, err := svc.persist.FindPermissionByID(id)
			if err != nil {
				return nil, status.New(codes.InvalidArgument, fmt.Sprintf("Permission: %s doesn't exist", id)).Err()
			}

			permissions = append(permissions, p)
		}
	}

	return permissions, nil
}
