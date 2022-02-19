package role_service

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
	role_proto "github.com/xdorro/golang-grpc-base-project/api/proto/role"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/handler"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService)

var _ role_proto.RoleServiceServer = (*RoleService)(nil)

type RoleService struct {
	log     *zap.Logger
	repo    *repo.Repo
	handler *handler.Handler
	redis   redis.UniversalClient

	// implement RoleService
	role_proto.UnimplementedRoleServiceServer
}

// NewService returns a new service instance
func NewService(
	log *zap.Logger, repo *repo.Repo, redis redis.UniversalClient, handler *handler.Handler, grpc *grpc.Server,
) role_proto.RoleServiceServer {
	svc := &RoleService{
		log:     log,
		repo:    repo,
		redis:   redis,
		handler: handler,
	}

	// Register RoleService Server
	role_proto.RegisterRoleServiceServer(grpc, svc)

	return svc
}

// FindAllRoles returns all roles
func (svc *RoleService) FindAllRoles(ctx context.Context, _ *role_proto.FindAllRolesRequest) (
	*role_proto.ListRolesResponse, error,
) {
	result := &role_proto.ListRolesResponse{}

	// Get value from redis
	if value := svc.redis.Get(ctx, common.FindAllRoles).Val(); value != "" {
		if err := protojson.Unmarshal([]byte(value), result); err != nil {
			return nil, err
		}
		return result, nil
	}

	roles := svc.repo.FindAllRoles()
	result = &role_proto.ListRolesResponse{Data: common.RolesProto(roles)}

	// Cache data to redis
	data, err := protojson.Marshal(result)
	if err != nil {
		return nil, err
	}
	svc.redis.Set(ctx, common.FindAllRoles, data, 5*time.Hour)

	return result, nil
}

// FindRoleByID returns a role by id
func (svc *RoleService) FindRoleByID(ctx context.Context, in *common_proto.UUIDRequest) (
	*role_proto.Role, error,
) {
	// Validate request
	if err := svc.handler.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	result := &role_proto.Role{}
	id := cast.ToUint64(in.GetId())

	// Get value from redis
	if value := svc.redis.Get(ctx, fmt.Sprintf(common.FindRoleByID, id)).Val(); value != "" {
		if err := protojson.Unmarshal([]byte(value), result); err != nil {
			return nil, err
		}
		return result, nil
	}

	r, err := svc.repo.FindRoleByID(id)
	if err != nil {
		return nil, common.RoleNotExist.Err()
	}

	result = common.RoleProto(r)

	// Cache data to redis
	data, err := protojson.Marshal(result)
	if err != nil {
		return nil, err
	}
	svc.redis.Set(ctx, fmt.Sprintf(common.FindRoleByID, id), data, 5*time.Hour)

	return result, nil
}

// CreateRole creates a new role
func (svc *RoleService) CreateRole(_ context.Context, in *role_proto.CreateRoleRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.handler.ValidateCreateRoleRequest(in); err != nil {
		svc.log.Error("svc.validateCreateRoleRequest()", zap.Error(err))
		return nil, err
	}

	slug := common.GetSlugOrMakeSlug(in.GetName(), in.GetSlug())
	if exist := svc.repo.ExistRoleBySlug(slug); exist {
		return nil, common.RoleAlreadyExist.Err()
	}

	permissions, err := svc.handler.ValidateListPermissions(in.GetPermissions())
	if err != nil {
		svc.log.Error("svc.validateCreateRolePermissions()", zap.Error(err))
		return nil, err
	}

	roleIn := &ent.Role{
		Name:       in.GetName(),
		Slug:       slug,
		Status:     in.GetStatus(),
		FullAccess: in.GetFullAccess(),
	}

	if err = svc.repo.CreateRole(roleIn, permissions); err != nil {
		return nil, err
	}

	if err = svc.redis.Del(svc.redis.Context(), common.KeyServiceRoles).Err(); err != nil {
		svc.log.Error("redis.Del()", zap.Error(err))
	}

	return common.Success.Proto(), nil
}

// UpdateRole update a role
func (svc *RoleService) UpdateRole(_ context.Context, in *role_proto.UpdateRoleRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.handler.ValidateUpdateRoleRequest(in); err != nil {
		svc.log.Error("svc.validateUpdateRoleRequest()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())
	r, err := svc.repo.FindRoleByID(id)
	if err != nil {
		return nil, common.RoleNotExist.Err()
	}

	slug := common.GetSlugOrMakeSlug(in.GetName(), in.GetSlug())
	if exist := svc.repo.ExistRoleByIDNotSlug(id, slug); exist {
		return nil, common.RoleAlreadyExist.Err()
	}

	permissions, err := svc.handler.ValidateListPermissions(in.GetPermissions())
	if err != nil {
		svc.log.Error("svc.validateUpdateRolePermissions()", zap.Error(err))
		return nil, err
	}

	r.Name = in.GetName()
	r.Slug = slug
	r.Status = in.GetStatus()
	r.FullAccess = in.GetFullAccess()

	if err = svc.repo.UpdateRole(r, permissions); err != nil {
		return nil, err
	}

	if err = svc.redis.Del(svc.redis.Context(), common.KeyServiceRoles).Err(); err != nil {
		svc.log.Error("redis.Del()", zap.Error(err))
	}

	return common.Success.Proto(), nil
}

// DeleteRole delete a role
func (svc *RoleService) DeleteRole(_ context.Context, in *common_proto.UUIDRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.handler.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())
	if exist := svc.repo.ExistRoleByID(id); !exist {
		return nil, common.RoleNotExist.Err()
	}

	if err := svc.repo.SoftDeleteRole(id); err != nil {
		return nil, err
	}

	if err := svc.redis.Del(svc.redis.Context(), common.KeyServiceRoles).Err(); err != nil {
		svc.log.Error("redis.Del()", zap.Error(err))
	}

	return common.Success.Proto(), nil
}
