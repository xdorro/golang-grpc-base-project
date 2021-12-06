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
	"github.com/kucow/golang-grpc-base-project/internal/validator"
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	authproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/auth"
	permissionproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/permission"
	roleproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/role"
	userproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/user"
)

type Service struct {
	log       *zap.Logger
	client    *ent.Client
	persist   *repo.Repo
	validator *validator.Validator
}

func NewService(opts *common.Option, srv *grpc.Server, validator *validator.Validator, persist *repo.Repo) {

	svc := &Service{
		log:       opts.Log,
		client:    opts.Client,
		persist:   persist,
		validator: validator,
	}

	// register Service Servers
	svc.registerServiceServers(opts, srv)

	// get Service Info
	svc.getServiceInfo(srv)
}

func (svc *Service) registerServiceServers(opts *common.Option, srv *grpc.Server) {
	// Register AuthService Server
	authproto.RegisterAuthServiceServer(srv, authservice.NewAuthService(opts, svc.validator, svc.persist))
	// Register UserService Server
	userproto.RegisterUserServiceServer(srv, userservice.NewUserService(opts, svc.validator, svc.persist))
	// Register RoleService Server
	roleproto.RegisterRoleServiceServer(srv, authservice.NewRoleService(opts, svc.validator, svc.persist))
	// Register PermissionService Server
	permissionproto.RegisterPermissionServiceServer(srv, authservice.NewPermissionService(opts, svc.validator, svc.persist))
}

func (svc *Service) getServiceInfo(srv *grpc.Server) {
	if viper.GetBool("SEEDER_SERVICE") {
		bulk := make([]*ent.PermissionCreate, 0)

		for name, val := range srv.GetServiceInfo() {
			if len(val.Methods) == 0 {
				return
			}

			for _, info := range val.Methods {
				slug := fmt.Sprintf("/%s/%s", name, info.Name)

				if !svc.persist.ExistPermissionBySlug(slug) {
					svc.log.Info("GetServiceInfo",
						zap.Any("Name", info.Name),
						zap.Any("Slug", slug),
					)

					bulk = append(bulk, svc.client.Permission.
						Create().
						SetName(info.Name).
						SetSlug(slug).
						SetStatus(1),
					)
				}
			}
		}

		if len(bulk) > 0 {
			if err := svc.persist.CreatePermissionBulk(bulk); err != nil {
				svc.log.Error("persist.CreatePermissionBulk()", zap.Error(err))
			}
		}
	}
}
