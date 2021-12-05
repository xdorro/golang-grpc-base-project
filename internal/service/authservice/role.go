package authservice

import (
	"context"

	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/kucow/golang-grpc-base-project/internal/common"
	"github.com/kucow/golang-grpc-base-project/internal/repo"
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	commonproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/common"
	roleproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/role"
	"github.com/kucow/golang-grpc-base-project/pkg/validator"
)

type RoleService struct {
	roleproto.UnimplementedRoleServiceServer

	log       *zap.Logger
	persist   repo.Persist
	validator *validator.Validator
}

func NewRoleService(opts *common.Option, persist repo.Persist) *RoleService {
	return &RoleService{
		log:       opts.Log,
		persist:   persist,
		validator: opts.Validator,
	}
}

func (svc *RoleService) FindAllRoles(context.Context, *roleproto.FindAllRolesRequest) (
	*roleproto.ListRolesResponse, error,
) {
	roles := svc.persist.FindAllRoles()

	return &roleproto.ListRolesResponse{Data: common.RolesProto(roles)}, nil
}

func (svc *RoleService) FindRoleByID(_ context.Context, in *commonproto.UUIDRequest) (
	*roleproto.Role, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.persist.FindRoleByID(in.Id)
	if err != nil {
		return nil, common.RoleNotExist.Err()
	}

	return common.RoleProto(u), nil
}

func (svc *RoleService) CreateRole(_ context.Context, in *roleproto.CreateRoleRequest) (
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
		Name:   in.GetName(),
		Slug:   in.GetSlug(),
		Status: in.GetStatus(),
	}

	if err = svc.persist.CreateRole(roleIn, permissions); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

func (svc *RoleService) UpdateRole(_ context.Context, in *roleproto.UpdateRoleRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateUpdateRoleRequest(in); err != nil {
		svc.log.Error("svc.validateUpdateRoleRequest()", zap.Error(err))
		return nil, err
	}

	r, err := svc.persist.FindRoleByID(in.GetId())
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

	if exist := svc.persist.ExistRoleByID(in.GetId()); !exist {
		return nil, common.RoleNotExist.Err()
	}

	if err := svc.persist.SoftDeleteRole(in.GetId()); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}
