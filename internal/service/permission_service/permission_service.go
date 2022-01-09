package permission_service

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	status_proto "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	common_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/common"
	permission_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/permission"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
	"github.com/xdorro/golang-grpc-base-project/internal/validator"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService)

type PermissionService struct {
	log       *zap.Logger
	repo      *repo.Repo
	validator *validator.Validator
	redis     redis.UniversalClient

	// implement PermissionService
	permission_proto.UnimplementedPermissionServiceServer
}

// NewService returns a new service instance
func NewService(
	log *zap.Logger, repo *repo.Repo, redis redis.UniversalClient, validator *validator.Validator, grpc *grpc.Server,
) *PermissionService {
	svc := &PermissionService{
		log:       log,
		repo:      repo,
		redis:     redis,
		validator: validator,
	}

	// Register PermissionService Server
	permission_proto.RegisterPermissionServiceServer(grpc, svc)

	return svc
}

// FindAllPermissions returns all permissions
func (svc *PermissionService) FindAllPermissions(context.Context, *permission_proto.FindAllPermissionsRequest) (
	*permission_proto.ListPermissionsResponse, error,
) {
	permissions := svc.repo.FindAllPermissions()

	return &permission_proto.ListPermissionsResponse{Data: common.PermissionsProto(permissions)}, nil
}

// FindPermissionByID returns a permission by id
func (svc *PermissionService) FindPermissionByID(_ context.Context, in *common_proto.UUIDRequest) (
	*permission_proto.Permission, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.repo.FindPermissionByID(cast.ToUint64(in.Id))
	if err != nil {
		return nil, common.PermissionNotExist.Err()
	}

	return common.PermissionProto(u), nil
}

// CreatePermission creates a new permission
func (svc *PermissionService) CreatePermission(_ context.Context, in *permission_proto.CreatePermissionRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCreatePermissionRequest(in); err != nil {
		svc.log.Error("svc.validateCreatePermissionRequest()", zap.Error(err))
		return nil, err
	}

	if exist := svc.repo.ExistPermissionBySlug(in.GetSlug()); exist {
		return nil, common.PermissionAlreadyExist.Err()
	}

	permissionIn := &ent.Permission{
		Name:   in.GetName(),
		Slug:   in.GetSlug(),
		Status: in.GetStatus(),
	}

	if err := svc.repo.CreatePermission(permissionIn); err != nil {
		return nil, err
	}

	if err := svc.redis.Del(svc.redis.Context(), common.KeyServiceRoles).Err(); err != nil {
		svc.log.Error("redis.Del()", zap.Error(err))
	}

	return common.Success.Proto(), nil
}

// UpdatePermission updates a permission
func (svc *PermissionService) UpdatePermission(_ context.Context, in *permission_proto.UpdatePermissionRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateUpdatePermissionRequest(in); err != nil {
		svc.log.Error("svc.validateUpdatePermissionRequest()", zap.Error(err))
		return nil, err
	}

	r, err := svc.repo.FindPermissionByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.PermissionNotExist.Err()
	}

	r.Name = in.GetName()
	r.Slug = in.GetSlug()
	r.Status = in.GetStatus()

	if err = svc.repo.UpdatePermission(r); err != nil {
		return nil, err
	}

	if err = svc.redis.Del(svc.redis.Context(), common.KeyServiceRoles).Err(); err != nil {
		svc.log.Error("redis.Del()", zap.Error(err))
	}

	return common.Success.Proto(), nil
}

// DeletePermission deletes a permission
func (svc *PermissionService) DeletePermission(_ context.Context, in *common_proto.UUIDRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())
	if exist := svc.repo.ExistPermissionByID(id); !exist {
		return nil, common.PermissionNotExist.Err()
	}

	if err := svc.repo.SoftDeletePermission(id); err != nil {
		return nil, err
	}

	if err := svc.redis.Del(svc.redis.Context(), common.KeyServiceRoles).Err(); err != nil {
		svc.log.Error("redis.Del()", zap.Error(err))
	}

	return common.Success.Proto(), nil
}
