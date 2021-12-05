package service

import (
	"google.golang.org/grpc"

	"github.com/kucow/golang-grpc-base-project/internal/common"
	"github.com/kucow/golang-grpc-base-project/internal/repo"
	"github.com/kucow/golang-grpc-base-project/internal/service/authservice"
	"github.com/kucow/golang-grpc-base-project/internal/service/permissionservice"
	"github.com/kucow/golang-grpc-base-project/internal/service/roleservice"
	"github.com/kucow/golang-grpc-base-project/internal/service/userservice"
	authproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/auth"
	permissionproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/permission"
	roleproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/role"
	userproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/user"
)

func NewService(opts *common.Option, srv *grpc.Server) {
	// Create new persist
	persist := repo.NewRepo(opts)

	// Register AuthService Server
	authproto.RegisterAuthServiceServer(srv, authservice.NewAuthService(opts, persist))
	// Register UserService Server
	userproto.RegisterUserServiceServer(srv, userservice.NewUserService(opts, persist))
	// Register RoleService Server
	roleproto.RegisterRoleServiceServer(srv, roleservice.NewRoleService(opts, persist))
	// Register PermissionService Server
	permissionproto.RegisterPermissionServiceServer(srv, permissionservice.NewPermissionService(opts, persist))
}
