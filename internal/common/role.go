package common

import (
	"github.com/spf13/cast"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/proto/role"
)

const (
	FindAllRoles = "dashboard:findAllRoles"
	FindRoleByID = "dashboard:%d:findRoleByID"
)

// RoleProto convert ent role to proto
func RoleProto(role *ent.Role) *role_proto.Role {
	return &role_proto.Role{
		Id:         cast.ToString(role.ID),
		Name:       role.Name,
		Slug:       role.Slug,
		FullAccess: role.FullAccess,
		Status:     role.Status,
	}
}

// RolesProto convert ent roles to proto
func RolesProto(roles []*ent.Role) []*role_proto.Role {
	result := make([]*role_proto.Role, len(roles))

	for index, role := range roles {
		result[index] = RoleProto(role)
	}

	return result
}
