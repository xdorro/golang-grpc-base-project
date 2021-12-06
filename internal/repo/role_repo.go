package repo

import (
	"time"

	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/pkg/ent"
	"github.com/xdorro/golang-grpc-base-project/pkg/ent/permission"
	"github.com/xdorro/golang-grpc-base-project/pkg/ent/role"
)

// FindAllRoles find all roles
func (repo *Repo) FindAllRoles() []*ent.Role {
	roles, err := repo.Client.Role.
		Query().
		Select(
			role.FieldID,
			role.FieldName,
			role.FieldSlug,
			role.FieldStatus,
		).
		Where(role.DeleteTimeIsNil()).
		All(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindAllRoles()", zap.Error(err))
		return nil
	}

	return roles
}

// FindRoleByID find role by ID
func (repo *Repo) FindRoleByID(id uint64) (*ent.Role, error) {
	r, err := repo.Client.Role.
		Query().
		Where(role.ID(id), role.DeleteTimeIsNil()).
		First(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindRoleByID()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindRoleBySlug find role by slug
func (repo *Repo) FindRoleBySlug(slug string) (*ent.Role, error) {
	r, err := repo.Client.Role.
		Query().
		Where(role.SlugEqualFold(slug), role.DeleteTimeIsNil()).
		First(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindRoleBySlug()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindRoleByIDAndPermissionID find role by ID and permissionID
func (repo *Repo) FindRoleByIDAndPermissionID(id, permissionId uint64) (*ent.Role, error) {
	r, err := repo.Client.Role.
		Query().
		Where(
			role.ID(id),
			role.DeleteTimeIsNil(),
			role.HasPermissionsWith(permission.ID(permissionId)),
		).
		First(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindRoleByIDAndPermissionID()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindRoleByIDAndPermissionIDNot find role by ID and permissionID not
func (repo *Repo) FindRoleByIDAndPermissionIDNot(id, permissionId uint64) (*ent.Role, error) {
	r, err := repo.Client.Role.
		Query().
		Where(
			role.ID(id),
			role.DeleteTimeIsNil(),
			role.HasPermissionsWith(permission.Not(permission.ID(permissionId))),
		).
		First(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindRoleByIDAndPermissionIDNot()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// ExistRoleByID exist role by ID
func (repo *Repo) ExistRoleByID(id uint64) bool {
	exist, err := repo.Client.Role.
		Query().
		Where(role.ID(id), role.DeleteTimeIsNil()).
		Exist(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.ExistRoleByID()", zap.Error(err))
		return exist
	}

	return exist
}

// ExistRoleBySlug exist role by slug
func (repo *Repo) ExistRoleBySlug(slug string) bool {
	exist, err := repo.Client.Role.
		Query().
		Where(role.Slug(slug), role.DeleteTimeIsNil()).
		Exist(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.ExistRoleBySlug()", zap.Error(err))
		return exist
	}

	return exist
}

// CreateRole create role
func (repo *Repo) CreateRole(r *ent.Role, p []*ent.Permission) error {
	r, err := repo.Client.Role.
		Create().
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		SetFullAccess(r.FullAccess).
		AddPermissions(p...).
		Save(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.CreateRole()", zap.Error(err))
		return err
	}

	return nil
}

// UpdateRole update role
func (repo *Repo) UpdateRole(r *ent.Role, p []*ent.Permission) error {
	_, err := repo.Client.Role.
		Update().
		Where(role.ID(r.ID), role.DeleteTimeIsNil()).
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		SetFullAccess(r.FullAccess).
		ClearPermissions().
		AddPermissions(p...).
		Save(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.UpdateRole()", zap.Error(err))
		return err
	}

	return nil
}

// SoftDeleteRole soft deletes role
func (repo *Repo) SoftDeleteRole(id uint64) error {
	_, err := repo.Client.Role.
		Update().
		Where(role.ID(id), role.DeleteTimeIsNil()).
		SetDeleteTime(time.Now()).
		ClearPermissions().
		Save(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.SoftDeleteRole()", zap.Error(err))
		return err
	}

	return nil
}
