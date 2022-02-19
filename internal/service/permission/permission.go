package permission_service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	status_proto "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/proto/common"
	permission_proto "github.com/xdorro/golang-grpc-base-project/api/proto/permission"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/handler"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService)

var _ permission_proto.PermissionServiceServer = (*PermissionService)(nil)

type PermissionService struct {
	log     *zap.Logger
	repo    *repo.Repo
	handler *handler.Handler
	redis   redis.UniversalClient

	// implement PermissionService
	permission_proto.UnimplementedPermissionServiceServer
}

// NewService returns a new service instance
func NewService(
	log *zap.Logger, repo *repo.Repo, redis redis.UniversalClient, handler *handler.Handler, grpc *grpc.Server,
) permission_proto.PermissionServiceServer {
	svc := &PermissionService{
		log:     log,
		repo:    repo,
		redis:   redis,
		handler: handler,
	}

	// Register PermissionService Server
	permission_proto.RegisterPermissionServiceServer(grpc, svc)

	return svc
}

// FindAllPermissions returns all permissions
func (svc *PermissionService) FindAllPermissions(ctx context.Context, _ *permission_proto.FindAllPermissionsRequest) (
	*permission_proto.ListPermissionsResponse, error,
) {
	result := &permission_proto.ListPermissionsResponse{}

	// Get value from redis
	if value := svc.redis.Get(ctx, common.FindAllPermissions).Val(); value != "" {
		if err := protojson.Unmarshal([]byte(value), result); err != nil {
			return nil, err
		}
		return result, nil
	}

	permissions := svc.repo.FindAllPermissions()
	result = &permission_proto.ListPermissionsResponse{Data: common.PermissionsProto(permissions)}

	// Cache data to redis
	data, err := protojson.Marshal(result)
	if err != nil {
		return nil, err
	}
	svc.redis.Set(ctx, common.FindAllPermissions, data, 5*time.Hour)

	return result, nil
}

// FindPermissionByID returns a permission by id
func (svc *PermissionService) FindPermissionByID(ctx context.Context, in *common_proto.UUIDRequest) (
	*permission_proto.Permission, error,
) {
	// Validate request
	if err := svc.handler.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	result := &permission_proto.Permission{}
	id := cast.ToUint64(in.GetId())

	// Get value from redis
	if value := svc.redis.Get(ctx, fmt.Sprintf(common.FindPermissionByID, id)).Val(); value != "" {
		if err := protojson.Unmarshal([]byte(value), result); err != nil {
			return nil, err
		}
		return result, nil
	}

	p, err := svc.repo.FindPermissionByID(id)
	if err != nil {
		return nil, common.PermissionNotExist.Err()
	}

	result = common.PermissionProto(p)

	// Cache data to redis
	data, err := protojson.Marshal(result)
	if err != nil {
		return nil, err
	}
	svc.redis.Set(ctx, fmt.Sprintf(common.FindPermissionByID, id), data, 5*time.Hour)

	return result, nil
}

// CreatePermission creates a new permission
func (svc *PermissionService) CreatePermission(_ context.Context, in *permission_proto.CreatePermissionRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.handler.ValidateCreatePermissionRequest(in); err != nil {
		svc.log.Error("svc.validateCreatePermissionRequest()", zap.Error(err))
		return nil, err
	}

	slug := common.GetSlugOrMakeSlug(in.GetName(), in.GetSlug())
	if exist := svc.repo.ExistPermissionBySlug(slug); exist {
		return nil, common.PermissionAlreadyExist.Err()
	}

	permissionIn := &ent.Permission{
		Name:   in.GetName(),
		Slug:   slug,
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
	if err := svc.handler.ValidateUpdatePermissionRequest(in); err != nil {
		svc.log.Error("svc.validateUpdatePermissionRequest()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())
	r, err := svc.repo.FindPermissionByID(id)
	if err != nil {
		return nil, common.PermissionNotExist.Err()
	}

	slug := common.GetSlugOrMakeSlug(in.GetName(), in.GetSlug())
	if exist := svc.repo.ExistPermissionByIDNotAndSlug(id, slug); exist {
		return nil, common.PermissionAlreadyExist.Err()
	}

	r.Name = in.GetName()
	r.Slug = slug
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
	if err := svc.handler.ValidateCommonID(in); err != nil {
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
