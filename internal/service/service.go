package service

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/event"
	"github.com/xdorro/golang-grpc-base-project/internal/validator"
	"github.com/xdorro/golang-grpc-base-project/pkg/client"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
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
	client     *client.Client
	grpcServer *grpc.Server

	validator *validator.Validator
	event     *event.Event
}

// NewService returns a new service instance
func NewService(
	client *client.Client,
	grpcServer *grpc.Server, redis redis.UniversalClient,
) *Service {
	svc := &Service{
		client:     client,
		grpcServer: grpcServer,
		redis:      redis,
	}

	// Create new validator
	svc.validator = validator.NewValidator(client)

	if viper.GetBool("ASYNQ_ENABLE") {
		// Create new event
		svc.event = event.NewEvent(client)
	}

	// register Service Servers
	svc.registerServiceServers()

	// get Service Info
	svc.getServiceInfo()

	return svc
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
					logger.Info("GetServiceInfo",
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
				logger.Error("persist.CreatePermissionBulk()", zap.Error(err))
			}
		}
	}
}

// Close closes the service.
func (svc *Service) Close() error {
	if svc.event != nil {
		return svc.event.Close()
	}

	return nil
}
