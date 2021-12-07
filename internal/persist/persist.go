package persist

import (
	"github.com/xdorro/golang-grpc-base-project/ent"
)

// Persist interface
type Persist interface {
	FindAllUsers() []*ent.User
	FindUserByEmail(email string) (*ent.User, error)
	FindUserByID(id uint64) (*ent.User, error)
	ExistUserByID(id uint64) bool
	ExistUserByEmail(email string) bool
	CreateUser(user *ent.User, roles []*ent.Role) error
	UpdateUser(user *ent.User, roles []*ent.Role) error
	DeleteUser(id uint64) error
	SoftDeleteUser(id uint64) error

	FindAllRoles() []*ent.Role
	FindRoleByID(id uint64) (*ent.Role, error)
	FindRoleBySlug(slug string) (*ent.Role, error)
	FindRoleByIDAndPermissionID(id, permissionId uint64) (*ent.Role, error)
	FindRoleByIDAndPermissionIDNot(id, permissionId uint64) (*ent.Role, error)
	ExistRoleByID(id uint64) bool
	ExistRoleBySlug(slug string) bool
	CreateRole(role *ent.Role, permissions []*ent.Permission) error
	UpdateRole(role *ent.Role, permissions []*ent.Permission) error
	SoftDeleteRole(id uint64) error

	FindAllPermissions() []*ent.Permission
	FindPermissionByID(id uint64) (*ent.Permission, error)
	FindPermissionBySlug(slug string) (*ent.Permission, error)
	FindPermissionByIDAndRoleIDNot(id uint64, roleId uint64) (*ent.Permission, error)
	ExistPermissionByID(id uint64) bool
	ExistPermissionBySlug(slug string) bool
	CreatePermission(role *ent.Permission) error
	CreatePermissionBulk(roles []*ent.PermissionCreate) error
	UpdatePermission(role *ent.Permission) error
	SoftDeletePermission(id uint64) error
}
