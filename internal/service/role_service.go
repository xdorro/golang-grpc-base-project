package service

import (
	"context"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/common"
	roleproto "github.com/xdorro/golang-grpc-base-project/proto/v1/role"
)

// FindAllRoles returns all roles
func (svc *Service) FindAllRoles(context.Context, *roleproto.FindAllRolesRequest) (
	*roleproto.ListRolesResponse, error,
) {
	roles := svc.client.Persist.FindAllRoles()

	return &roleproto.ListRolesResponse{Data: common.RolesProto(roles)}, nil
}

// FindRoleByID returns a role by id
func (svc *Service) FindRoleByID(_ context.Context, in *commonproto.UUIDRequest) (
	*roleproto.Role, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.client.Persist.FindRoleByID(cast.ToUint64(in.Id))
	if err != nil {
		return nil, common.RoleNotExist.Err()
	}

	return common.RoleProto(u), nil
}

// CreateRole creates a new role
func (svc *Service) CreateRole(_ context.Context, in *roleproto.CreateRoleRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCreateRoleRequest(in); err != nil {
		svc.log.Error("svc.validateCreateRoleRequest()", zap.Error(err))
		return nil, err
	}

	if exist := svc.client.Persist.ExistRoleBySlug(in.GetSlug()); exist {
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

	if err = svc.client.Persist.CreateRole(roleIn, permissions); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// UpdateRole update a role
func (svc *Service) UpdateRole(_ context.Context, in *roleproto.UpdateRoleRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateUpdateRoleRequest(in); err != nil {
		svc.log.Error("svc.validateUpdateRoleRequest()", zap.Error(err))
		return nil, err
	}

	r, err := svc.client.Persist.FindRoleByID(cast.ToUint64(in.GetId()))
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

	if err = svc.client.Persist.UpdateRole(r, permissions); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// DeleteRole delete a role
func (svc *Service) DeleteRole(_ context.Context, in *commonproto.UUIDRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())
	if exist := svc.client.Persist.ExistRoleByID(id); !exist {
		return nil, common.RoleNotExist.Err()
	}

	if err := svc.client.Persist.SoftDeleteRole(id); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}
