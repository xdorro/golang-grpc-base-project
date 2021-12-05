package common

import (
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	roleproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/role"
)

func RoleProto(role *ent.Role) *roleproto.Role {
	return &roleproto.Role{
		Id:     role.ID,
		Name:   role.Name,
		Slug:   role.Slug,
		Status: role.Status,
	}
}

func RolesProto(roles []*ent.Role) []*roleproto.Role {
	result := make([]*roleproto.Role, len(roles))

	for index, role := range roles {
		result[index] = RoleProto(role)
	}

	return result
}
