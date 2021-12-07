package authservice

import (
	"context"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/common/optioncommon"
	"github.com/xdorro/golang-grpc-base-project/internal/persist"
	"github.com/xdorro/golang-grpc-base-project/pkg/validator"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/common"
	permissionproto2 "github.com/xdorro/golang-grpc-base-project/proto/v1/permission"
)

type PermissionService struct {
	permissionproto2.UnimplementedPermissionServiceServer

	log       *zap.Logger
	persist   persist.Persist
	validator *validator.Validator
}

func NewPermissionService(
	opts *optioncommon.Option, validator *validator.Validator, persist persist.Persist,
) *PermissionService {
	svc := &PermissionService{
		log:       opts.Log,
		persist:   persist,
		validator: validator,
	}

	return svc
}

func (svc *PermissionService) FindAllPermissions(context.Context, *permissionproto2.FindAllPermissionsRequest) (
	*permissionproto2.ListPermissionsResponse, error,
) {
	permissions := svc.persist.FindAllPermissions()

	return &permissionproto2.ListPermissionsResponse{Data: common.PermissionsProto(permissions)}, nil
}

func (svc *PermissionService) FindPermissionByID(_ context.Context, in *commonproto.UUIDRequest) (
	*permissionproto2.Permission, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.persist.FindPermissionByID(cast.ToUint64(in.Id))
	if err != nil {
		return nil, common.PermissionNotExist.Err()
	}

	return common.PermissionProto(u), nil
}

func (svc *PermissionService) CreatePermission(_ context.Context, in *permissionproto2.CreatePermissionRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCreatePermissionRequest(in); err != nil {
		svc.log.Error("svc.validateCreatePermissionRequest()", zap.Error(err))
		return nil, err
	}

	if exist := svc.persist.ExistPermissionBySlug(in.GetSlug()); exist {
		return nil, common.PermissionAlreadyExist.Err()
	}

	permissionIn := &ent.Permission{
		Name:   in.GetName(),
		Slug:   in.GetSlug(),
		Status: in.GetStatus(),
	}

	if err := svc.persist.CreatePermission(permissionIn); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

func (svc *PermissionService) UpdatePermission(_ context.Context, in *permissionproto2.UpdatePermissionRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateUpdatePermissionRequest(in); err != nil {
		svc.log.Error("svc.validateUpdatePermissionRequest()", zap.Error(err))
		return nil, err
	}

	r, err := svc.persist.FindPermissionByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.PermissionNotExist.Err()
	}

	r.Name = in.GetName()
	r.Slug = in.GetSlug()
	r.Status = in.GetStatus()

	if err = svc.persist.UpdatePermission(r); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

func (svc *PermissionService) DeletePermission(_ context.Context, in *commonproto.UUIDRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())
	if exist := svc.persist.ExistPermissionByID(id); !exist {
		return nil, common.PermissionNotExist.Err()
	}

	if err := svc.persist.SoftDeletePermission(id); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}
