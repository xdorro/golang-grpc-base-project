package service

import (
	"github.com/google/wire"

	auth_proto "github.com/xdorro/golang-grpc-base-project/api/proto/auth"
	permission_proto "github.com/xdorro/golang-grpc-base-project/api/proto/permission"
	role_proto "github.com/xdorro/golang-grpc-base-project/api/proto/role"
	user_proto "github.com/xdorro/golang-grpc-base-project/api/proto/user"
	auth_service "github.com/xdorro/golang-grpc-base-project/internal/service/auth"
	permission_service "github.com/xdorro/golang-grpc-base-project/internal/service/permission"
	role_service "github.com/xdorro/golang-grpc-base-project/internal/service/role"
	user_service "github.com/xdorro/golang-grpc-base-project/internal/service/user"
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
	user_proto.UserServiceServer
	auth_proto.AuthServiceServer
	role_proto.RoleServiceServer
	permission_proto.PermissionServiceServer
}
