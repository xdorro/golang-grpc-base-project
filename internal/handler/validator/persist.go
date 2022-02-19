package validator_handler

import (
	"github.com/xdorro/golang-grpc-base-project/api/ent"
	auth_proto "github.com/xdorro/golang-grpc-base-project/api/proto/auth"
	common_proto "github.com/xdorro/golang-grpc-base-project/api/proto/common"
	permission_proto "github.com/xdorro/golang-grpc-base-project/api/proto/permission"
	role_proto "github.com/xdorro/golang-grpc-base-project/api/proto/role"
	user_proto "github.com/xdorro/golang-grpc-base-project/api/proto/user"
)

type ValidatorPersist interface {
	ValidateError(err error) error
	ValidateCommonID(in *common_proto.UUIDRequest) error
	ValidateCommonSlug(in *common_proto.SlugRequest) error

	ValidateLoginRequest(in *auth_proto.LoginRequest) error
	ValidateTokenRequest(in *auth_proto.TokenRequest) error

	ValidateCreateUserRequest(in *user_proto.CreateUserRequest) error
	ValidateUpdateUserRequest(in *user_proto.UpdateUserRequest) error

	ValidateCreateRoleRequest(in *role_proto.CreateRoleRequest) error
	ValidateUpdateRoleRequest(in *role_proto.UpdateRoleRequest) error
	ValidateRole(slug string) (*ent.Role, error)

	ValidateCreatePermissionRequest(in *permission_proto.CreatePermissionRequest) error
	ValidateUpdatePermissionRequest(in *permission_proto.UpdatePermissionRequest) error
	ValidateListPermissions(list []string) ([]*ent.Permission, error)
}
