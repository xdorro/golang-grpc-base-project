package permission_repo

import (
	"github.com/xdorro/golang-grpc-base-project/api/ent"
)

type PermissionPersist interface {
	FindAllPermissions() []*ent.Permission
	FindPermissionByID(id uint64) (*ent.Permission, error)
	FindPermissionBySlug(slug string) (*ent.Permission, error)
	FindPermissionByIDAndRoleIDNot(id uint64, roleID uint64) (*ent.Permission, error)
	ExistPermissionByID(id uint64) bool
	ExistPermissionBySlug(slug string) bool
	CreatePermission(role *ent.Permission) error
	CreatePermissionBulk(roles []*ent.PermissionCreate) error
	UpdatePermission(role *ent.Permission) error
	SoftDeletePermission(id uint64) error
}
