package permissionservice

import (
	"github.com/kucow/golang-grpc-base/pkg/ent"
	permissionproto "github.com/kucow/golang-grpc-base/pkg/proto/v1/permission"
)

func PermissionProto(permission *ent.Permission) *permissionproto.Permission {
	return &permissionproto.Permission{
		Id:     permission.ID,
		Name:   permission.Name,
		Slug:   permission.Slug,
		Status: permission.Status,
	}
}

func PermissionsProto(permissions []*ent.Permission) []*permissionproto.Permission {
	result := make([]*permissionproto.Permission, len(permissions))

	for index, permission := range permissions {
		result[index] = PermissionProto(permission)
	}

	return result
}
