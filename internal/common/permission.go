package common

import (
	"github.com/spf13/cast"

	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	permissionproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/permission"
)

// PermissionProto convert ent permission to proto
func PermissionProto(permission *ent.Permission) *permissionproto.Permission {
	return &permissionproto.Permission{
		Id:     cast.ToString(permission.ID),
		Name:   permission.Name,
		Slug:   permission.Slug,
		Status: permission.Status,
	}
}

// PermissionsProto convert ent permissions to proto
func PermissionsProto(permissions []*ent.Permission) []*permissionproto.Permission {
	result := make([]*permissionproto.Permission, len(permissions))

	for index, permission := range permissions {
		result[index] = PermissionProto(permission)
	}

	return result
}
