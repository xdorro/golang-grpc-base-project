package service

import (
	"github.com/google/wire"

	"github.com/xdorro/golang-grpc-base-project/internal/service/auth_service"
	"github.com/xdorro/golang-grpc-base-project/internal/service/permission_service"
	"github.com/xdorro/golang-grpc-base-project/internal/service/role_service"
	"github.com/xdorro/golang-grpc-base-project/internal/service/user_service"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	wire.Struct(new(Service), "*"),
	user_service.ProviderSet,
	auth_service.ProviderSet,
	role_service.ProviderSet,
	permission_service.ProviderSet,
)

type Service struct {
	*user_service.UserService
	*auth_service.AuthService
	*role_service.RoleService
	*permission_service.PermissionService
}
