package role_repo

import (
	"github.com/xdorro/golang-grpc-base-project/api/ent"
)

type RolePersist interface {
	FindAllRoles() []*ent.Role

	FindRoleByID(id uint64) (*ent.Role, error)
	FindRoleBySlug(slug string) (*ent.Role, error)
	FindRoleByIDAndPermissionID(id, permissionID uint64) (*ent.Role, error)
	FindRoleByIDAndPermissionIDNot(id, permissionID uint64) (*ent.Role, error)

	ExistRoleByID(id uint64) bool
	ExistRoleBySlug(slug string) bool
	ExistRoleByIDNotSlug(id uint64, slug string) bool

	CreateRole(role *ent.Role, permissions []*ent.Permission) error

	UpdateRole(role *ent.Role, permissions []*ent.Permission) error

	SoftDeleteRole(id uint64) error
}
