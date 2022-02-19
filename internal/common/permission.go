package common

import (
	"github.com/spf13/cast"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/proto/permission"
)

const (
	FindAllPermissions = "dashboard:findAllPermissions"
	FindPermissionByID = "dashboard:%d:findPermissionByID"
)

// PermissionProto convert ent permission to proto
func PermissionProto(permission *ent.Permission) *permission_proto.Permission {
	return &permission_proto.Permission{
		Id:     cast.ToString(permission.ID),
		Name:   permission.Name,
		Slug:   permission.Slug,
		Status: permission.Status,
	}
}

// PermissionsProto convert ent permissions to proto
func PermissionsProto(permissions []*ent.Permission) []*permission_proto.Permission {
	result := make([]*permission_proto.Permission, len(permissions))

	for index, permission := range permissions {
		result[index] = PermissionProto(permission)
	}

	return result
}
