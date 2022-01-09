package role_service

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
	role_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/role"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
	"github.com/xdorro/golang-grpc-base-project/internal/validator"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService)

type RoleService struct {
	log       *zap.Logger
	repo      *repo.Repo
	validator *validator.Validator
	redis     redis.UniversalClient

	// implement RoleService
	role_proto.UnimplementedRoleServiceServer
}

// NewService returns a new service instance
func NewService(
	log *zap.Logger, repo *repo.Repo, redis redis.UniversalClient, validator *validator.Validator, grpc *grpc.Server,
) *RoleService {
	svc := &RoleService{
		log:       log,
		repo:      repo,
		redis:     redis,
		validator: validator,
	}

	// Register RoleService Server
	role_proto.RegisterRoleServiceServer(grpc, svc)

	return svc
}

// FindAllRoles returns all roles
func (svc *RoleService) FindAllRoles(context.Context, *role_proto.FindAllRolesRequest) (
	*role_proto.ListRolesResponse, error,
) {
	roles := svc.repo.FindAllRoles()

	return &role_proto.ListRolesResponse{Data: common.RolesProto(roles)}, nil
}

// FindRoleByID returns a role by id
func (svc *RoleService) FindRoleByID(_ context.Context, in *common_proto.UUIDRequest) (
	*role_proto.Role, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.repo.FindRoleByID(cast.ToUint64(in.Id))
	if err != nil {
		return nil, common.RoleNotExist.Err()
	}

	return common.RoleProto(u), nil
}

// CreateRole creates a new role
func (svc *RoleService) CreateRole(_ context.Context, in *role_proto.CreateRoleRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCreateRoleRequest(in); err != nil {
		svc.log.Error("svc.validateCreateRoleRequest()", zap.Error(err))
		return nil, err
	}

	if exist := svc.repo.ExistRoleBySlug(in.GetSlug()); exist {
		return nil, common.RoleAlreadyExist.Err()
	}

	permissions, err := svc.validator.ValidateListPermissions(in.GetPermissions())
	if err != nil {
		svc.log.Error("svc.validateCreateRolePermissions()", zap.Error(err))
		return nil, err
	}

	roleIn := &ent.Role{
		Name:       in.GetName(),
		Slug:       in.GetSlug(),
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
	if err := svc.validator.ValidateUpdateRoleRequest(in); err != nil {
		svc.log.Error("svc.validateUpdateRoleRequest()", zap.Error(err))
		return nil, err
	}

	r, err := svc.repo.FindRoleByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.RoleNotExist.Err()
	}

	permissions, err := svc.validator.ValidateListPermissions(in.GetPermissions())
	if err != nil {
		svc.log.Error("svc.validateUpdateRolePermissions()", zap.Error(err))
		return nil, err
	}

	r.Name = in.GetName()
	r.Slug = in.GetSlug()
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
	if err := svc.validator.ValidateCommonID(in); err != nil {
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
