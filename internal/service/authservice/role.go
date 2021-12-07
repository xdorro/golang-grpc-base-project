package authservice

import (
	"context"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
	"github.com/xdorro/golang-grpc-base-project/internal/validator"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/common"
	roleproto2 "github.com/xdorro/golang-grpc-base-project/proto/v1/role"
)

type RoleService struct {
	roleproto2.UnimplementedRoleServiceServer

	log       *zap.Logger
	persist   repo.Persist
	validator *validator.Validator
}

func NewRoleService(opts *common.Option, validator *validator.Validator, persist repo.Persist) *RoleService {
	return &RoleService{
		log:       opts.Log,
		persist:   persist,
		validator: validator,
	}
}

func (svc *RoleService) FindAllRoles(context.Context, *roleproto2.FindAllRolesRequest) (
	*roleproto2.ListRolesResponse, error,
) {
	roles := svc.persist.FindAllRoles()

	return &roleproto2.ListRolesResponse{Data: common.RolesProto(roles)}, nil
}

func (svc *RoleService) FindRoleByID(_ context.Context, in *commonproto.UUIDRequest) (
	*roleproto2.Role, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.persist.FindRoleByID(cast.ToUint64(in.Id))
	if err != nil {
		return nil, common.RoleNotExist.Err()
	}

	return common.RoleProto(u), nil
}

func (svc *RoleService) CreateRole(_ context.Context, in *roleproto2.CreateRoleRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCreateRoleRequest(in); err != nil {
		svc.log.Error("svc.validateCreateRoleRequest()", zap.Error(err))
		return nil, err
	}

	if exist := svc.persist.ExistRoleBySlug(in.GetSlug()); exist {
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

	if err = svc.persist.CreateRole(roleIn, permissions); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

func (svc *RoleService) UpdateRole(_ context.Context, in *roleproto2.UpdateRoleRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateUpdateRoleRequest(in); err != nil {
		svc.log.Error("svc.validateUpdateRoleRequest()", zap.Error(err))
		return nil, err
	}

	r, err := svc.persist.FindRoleByID(cast.ToUint64(in.GetId()))
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

	if err = svc.persist.UpdateRole(r, permissions); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

func (svc *RoleService) DeleteRole(_ context.Context, in *commonproto.UUIDRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())
	if exist := svc.persist.ExistRoleByID(id); !exist {
		return nil, common.RoleNotExist.Err()
	}

	if err := svc.persist.SoftDeleteRole(id); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}
