package permissionservice

import (
	"context"

	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/internal/repo"
	"github.com/kucow/golang-grpc-base/pkg/ent"
	commonproto "github.com/kucow/golang-grpc-base/pkg/proto/v1/common"
	permissionproto "github.com/kucow/golang-grpc-base/pkg/proto/v1/permission"
)

type PermissionService struct {
	permissionproto.UnimplementedPermissionServiceServer

	log     *zap.Logger
	persist repo.Persist
}

func NewPermissionService(opts *common.Option, persist repo.Persist) *PermissionService {
	return &PermissionService{
		log:     opts.Log,
		persist: persist,
	}
}

func (svc *PermissionService) FindAllPermissions(context.Context, *permissionproto.FindAllPermissionsRequest) (
	*permissionproto.ListPermissionsResponse, error,
) {
	permissions := svc.persist.FindAllPermissions()

	return &permissionproto.ListPermissionsResponse{Data: PermissionsProto(permissions)}, nil
}

func (svc *PermissionService) FindPermissionByID(_ context.Context, in *commonproto.UUIDRequest) (
	*permissionproto.Permission, error,
) {
	// Validate request
	if err := common.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.persist.FindPermissionByID(in.Id)
	if err != nil {
		return nil, common.PermissionNotExist.Err()
	}

	return PermissionProto(u), nil
}

func (svc *PermissionService) CreatePermission(_ context.Context, in *permissionproto.CreatePermissionRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validateCreatePermissionRequest(in); err != nil {
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
	if err := svc.validateUpdatePermissionRequest(in); err != nil {
		svc.log.Error("svc.validateUpdatePermissionRequest()", zap.Error(err))
		return nil, err
	}

	r, err := svc.persist.FindPermissionByID(in.GetId())
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
	if err := common.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	if exist := svc.persist.ExistPermissionByID(in.GetId()); !exist {
		return nil, common.PermissionNotExist.Err()
	}

	if err := svc.persist.SoftDeletePermission(in.GetId()); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}
