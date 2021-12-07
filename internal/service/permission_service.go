package service

import (
	"context"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/pkg/common"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/common"
	permissionproto "github.com/xdorro/golang-grpc-base-project/proto/v1/permission"
)

// FindAllPermissions returns all permissions
func (svc *Service) FindAllPermissions(context.Context, *permissionproto.FindAllPermissionsRequest) (
	*permissionproto.ListPermissionsResponse, error,
) {
	permissions := svc.client.Persist.FindAllPermissions()

	return &permissionproto.ListPermissionsResponse{Data: common.PermissionsProto(permissions)}, nil
}

// FindPermissionByID returns a permission by id
func (svc *Service) FindPermissionByID(_ context.Context, in *commonproto.UUIDRequest) (
	*permissionproto.Permission, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.client.Persist.FindPermissionByID(cast.ToUint64(in.Id))
	if err != nil {
		return nil, common.PermissionNotExist.Err()
	}

	return common.PermissionProto(u), nil
}

// CreatePermission creates a new permission
func (svc *Service) CreatePermission(_ context.Context, in *permissionproto.CreatePermissionRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCreatePermissionRequest(in); err != nil {
		svc.log.Error("svc.validateCreatePermissionRequest()", zap.Error(err))
		return nil, err
	}

	if exist := svc.client.Persist.ExistPermissionBySlug(in.GetSlug()); exist {
		return nil, common.PermissionAlreadyExist.Err()
	}

	permissionIn := &ent.Permission{
		Name:   in.GetName(),
		Slug:   in.GetSlug(),
		Status: in.GetStatus(),
	}

	if err := svc.client.Persist.CreatePermission(permissionIn); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// UpdatePermission updates a permission
func (svc *Service) UpdatePermission(_ context.Context, in *permissionproto.UpdatePermissionRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateUpdatePermissionRequest(in); err != nil {
		svc.log.Error("svc.validateUpdatePermissionRequest()", zap.Error(err))
		return nil, err
	}

	r, err := svc.client.Persist.FindPermissionByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.PermissionNotExist.Err()
	}

	r.Name = in.GetName()
	r.Slug = in.GetSlug()
	r.Status = in.GetStatus()

	if err = svc.client.Persist.UpdatePermission(r); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// DeletePermission deletes a permission
func (svc *Service) DeletePermission(_ context.Context, in *commonproto.UUIDRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())
	if exist := svc.client.Persist.ExistPermissionByID(id); !exist {
		return nil, common.PermissionNotExist.Err()
	}

	if err := svc.client.Persist.SoftDeletePermission(id); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}
