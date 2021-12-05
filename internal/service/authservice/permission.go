package authservice

import (
	"context"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/kucow/golang-grpc-base-project/internal/common"
	"github.com/kucow/golang-grpc-base-project/internal/repo"
	"github.com/kucow/golang-grpc-base-project/internal/validator"
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	commonproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/common"
	permissionproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/permission"
)

type PermissionService struct {
	permissionproto.UnimplementedPermissionServiceServer

	log       *zap.Logger
	persist   repo.Persist
	validator *validator.Validator
}

func NewPermissionService(opts *common.Option, persist repo.Persist) *PermissionService {
	svc := &PermissionService{
		log:       opts.Log,
		persist:   persist,
		validator: opts.Validator,
	}

	return svc
}

func (svc *PermissionService) FindAllPermissions(context.Context, *permissionproto.FindAllPermissionsRequest) (
	*permissionproto.ListPermissionsResponse, error,
) {
	permissions := svc.persist.FindAllPermissions()

	return &permissionproto.ListPermissionsResponse{Data: common.PermissionsProto(permissions)}, nil
}

func (svc *PermissionService) FindPermissionByID(_ context.Context, in *commonproto.UUIDRequest) (
	*permissionproto.Permission, error,
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

func (svc *PermissionService) CreatePermission(_ context.Context, in *permissionproto.CreatePermissionRequest) (
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

func (svc *PermissionService) UpdatePermission(_ context.Context, in *permissionproto.UpdatePermissionRequest) (
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
