package userservice

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/pkg/ent"
	userproto "github.com/kucow/golang-grpc-base/pkg/proto/v1/user"
)

func (svc *UserService) validateCreateUserRequest(in *userproto.CreateUserRequest) error {
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
			validation.Each(common.IsULID),
		),
	)

	return common.ValidateError(err)
}

func (svc *UserService) validateUpdateUserRequest(in *userproto.UpdateUserRequest) error {
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
			validation.Each(common.IsULID),
		),
	)

	return common.ValidateError(err)
}

func (svc *UserService) validateListRoles(list []string) ([]*ent.Role, error) {
	roles := make([]*ent.Role, 0)

	if len(list) > 0 {
		for _, id := range list {
			r, err := svc.persist.FindRoleByID(id)
			if err != nil {
				return nil, status.New(codes.InvalidArgument, fmt.Sprintf("role: %s doesn't exist", id)).Err()
			}

			roles = append(roles, r)
		}
	}

	return roles, nil
}
