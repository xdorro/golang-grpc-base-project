package service

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kucow/golang-grpc-base-project/internal/common"
	"github.com/kucow/golang-grpc-base-project/internal/repo"
	"github.com/kucow/golang-grpc-base-project/internal/service/authservice"
	"github.com/kucow/golang-grpc-base-project/internal/service/userservice"
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	authproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/auth"
	permissionproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/permission"
	roleproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/role"
	userproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/user"
	"github.com/kucow/golang-grpc-base-project/pkg/validator"
)

func NewService(opts *common.Option, srv *grpc.Server) {
	// Create new persist
	persist := repo.NewRepo(opts.Ctx, opts.Log, opts.Client)

	// Create new validator
	opts.Validator = validator.NewValidator(opts.Log, persist)

	// register Service Servers
	registerServiceServers(opts, srv, persist)

	// get Service Info
	getServiceInfo(opts, srv, persist)
}

func registerServiceServers(opts *common.Option, srv *grpc.Server, persist repo.Persist) {
	// Register AuthService Server
	authproto.RegisterAuthServiceServer(srv, authservice.NewAuthService(opts, persist))
	// Register UserService Server
	userproto.RegisterUserServiceServer(srv, userservice.NewUserService(opts, persist))
	// Register RoleService Server
	roleproto.RegisterRoleServiceServer(srv, authservice.NewRoleService(opts, persist))
	// Register PermissionService Server
	permissionproto.RegisterPermissionServiceServer(srv, authservice.NewPermissionService(opts, persist))
}

func getServiceInfo(opts *common.Option, srv *grpc.Server, persist repo.Persist) {
	if viper.GetBool("SEEDER_SERVICE") {
		bulk := make([]*ent.PermissionCreate, 0)

		for name, val := range srv.GetServiceInfo() {
			if len(val.Methods) == 0 {
				return
			}

			for _, info := range val.Methods {
				slug := fmt.Sprintf("%s/%s", name, info.Name)
				if !persist.ExistPermissionBySlug(slug) {
					opts.Log.Info("GetServiceInfo",
						zap.Any("Name", info.Name),
						zap.Any("Slug", slug),
					)

					bulk = append(bulk, opts.Client.Permission.
						Create().
						SetName(info.Name).
						SetSlug(slug).
						SetStatus(1),
					)
				}
			}
		}

		if len(bulk) > 0 {
			if err := persist.CreatePermissionBulk(bulk); err != nil {
				opts.Log.Error("persist.CreatePermissionBulk()", zap.Error(err))
			}
		}
	}
}
