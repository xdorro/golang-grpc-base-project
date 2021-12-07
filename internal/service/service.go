package service

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/validator"
	"github.com/xdorro/golang-grpc-base-project/pkg/client"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/auth"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/permission"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/role"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/user"
)

// Service is the service struct
type Service struct {
	authproto.UnimplementedAuthServiceServer
	permissionproto.UnimplementedPermissionServiceServer
	roleproto.UnimplementedRoleServiceServer
	userproto.UnimplementedUserServiceServer

	redis      redis.UniversalClient
	log        *zap.Logger
	client     *client.Client
	validator  *validator.Validator
	grpcServer *grpc.Server
}

// NewService returns a new service instance
func NewService(
	log *zap.Logger, client *client.Client, validator *validator.Validator,
	grpcServer *grpc.Server, redis redis.UniversalClient,
) {
	svc := &Service{
		log:        log,
		client:     client,
		validator:  validator,
		grpcServer: grpcServer,
		redis:      redis,
	}

	// register Service Servers
	svc.registerServiceServers()

	// get Service Info
	svc.getServiceInfo()
}

// registerServiceServers registers Service Servers
func (svc *Service) registerServiceServers() {
	// Register AuthService Server
	authproto.RegisterAuthServiceServer(svc.grpcServer, svc)
	// Register UserService Server
	userproto.RegisterUserServiceServer(svc.grpcServer, svc)
	// Register RoleService Server
	roleproto.RegisterRoleServiceServer(svc.grpcServer, svc)
	// Register PermissionService Server
	permissionproto.RegisterPermissionServiceServer(svc.grpcServer, svc)
}

// getServiceInfo returns service info
func (svc *Service) getServiceInfo() {
	if viper.GetBool("SEEDER_SERVICE") {
		bulk := make([]*ent.PermissionCreate, 0)

		for name, val := range svc.grpcServer.GetServiceInfo() {
			if len(val.Methods) == 0 {
				return
			}

			for _, info := range val.Methods {
				slug := fmt.Sprintf("/%s/%s", name, info.Name)

				if !svc.client.Persist.ExistPermissionBySlug(slug) {
					svc.log.Info("GetServiceInfo",
						zap.Any("Name", info.Name),
						zap.Any("Slug", slug),
					)

					bulk = append(bulk, svc.client.DB.Permission.
						Create().
						SetName(info.Name).
						SetSlug(slug).
						SetStatus(1),
					)
				}
			}
		}

		if len(bulk) > 0 {
			if err := svc.client.Persist.CreatePermissionBulk(bulk); err != nil {
				svc.log.Error("persist.CreatePermissionBulk()", zap.Error(err))
			}
		}
	}
}
