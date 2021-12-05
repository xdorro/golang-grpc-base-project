package repo

import (
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
)

// Persist interface
type Persist interface {
	FindAllUsers() []*ent.User
	FindUserByEmail(email string) (*ent.User, error)
	FindUserByID(id string) (*ent.User, error)
	ExistUserByID(id string) bool
	ExistUserByEmail(email string) bool
	CreateUser(user *ent.User, roles []*ent.Role) error
	UpdateUser(user *ent.User, roles []*ent.Role) error
	DeleteUser(id string) error
	SoftDeleteUser(id string) error

	FindAllRoles() []*ent.Role
	FindRoleByID(id string) (*ent.Role, error)
	FindRoleByIDAndPermissionID(id, permissionId string) (*ent.Role, error)
	FindRoleByIDAndPermissionIDNot(id, permissionId string) (*ent.Role, error)
	ExistRoleByID(id string) bool
	ExistRoleBySlug(slug string) bool
	CreateRole(role *ent.Role, permissions []*ent.Permission) error
	UpdateRole(role *ent.Role, permissions []*ent.Permission) error
	SoftDeleteRole(id string) error

	FindAllPermissions() []*ent.Permission
	FindPermissionByID(id string) (*ent.Permission, error)
	FindPermissionByIDAndRoleIDNot(id string, roleId string) (*ent.Permission, error)
	ExistPermissionByID(id string) bool
	ExistPermissionBySlug(slug string) bool
	CreatePermission(role *ent.Permission) error
	UpdatePermission(role *ent.Permission) error
	SoftDeletePermission(id string) error
}
