package repo

import (
	"time"

	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/ent/permission"
	"github.com/xdorro/golang-grpc-base-project/ent/role"
)

// FindAllRoles find all roles
func (repo *Repo) FindAllRoles() []*ent.Role {
	roles, err := repo.client.Role.
		Query().
		Select(
			role.FieldID,
			role.FieldName,
			role.FieldSlug,
			role.FieldStatus,
		).
		Where(role.DeleteTimeIsNil()).
		All(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindAllRoles()", zap.Error(err))
		return nil
	}

	return roles
}

// FindRoleByID find role by ID
func (repo *Repo) FindRoleByID(id uint64) (*ent.Role, error) {
	r, err := repo.client.Role.
		Query().
		Where(role.ID(id), role.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindRoleByID()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindRoleBySlug find role by slug
func (repo *Repo) FindRoleBySlug(slug string) (*ent.Role, error) {
	r, err := repo.client.Role.
		Query().
		Where(role.SlugEqualFold(slug), role.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindRoleBySlug()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindRoleByIDAndPermissionID find role by ID and permissionID
func (repo *Repo) FindRoleByIDAndPermissionID(id, permissionID uint64) (*ent.Role, error) {
	r, err := repo.client.Role.
		Query().
		Where(
			role.ID(id),
			role.DeleteTimeIsNil(),
			role.HasPermissionsWith(permission.ID(permissionID)),
		).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindRoleByIDAndPermissionID()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindRoleByIDAndPermissionIDNot find role by ID and permissionID not
func (repo *Repo) FindRoleByIDAndPermissionIDNot(id, permissionID uint64) (*ent.Role, error) {
	r, err := repo.client.Role.
		Query().
		Where(
			role.ID(id),
			role.DeleteTimeIsNil(),
			role.HasPermissionsWith(permission.Not(permission.ID(permissionID))),
		).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindRoleByIDAndPermissionIDNot()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// ExistRoleByID exist role by ID
func (repo *Repo) ExistRoleByID(id uint64) bool {
	exist, err := repo.client.Role.
		Query().
		Where(role.ID(id), role.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		repo.log.Error("persist.ExistRoleByID()", zap.Error(err))
		return exist
	}

	return exist
}

// ExistRoleBySlug exist role by slug
func (repo *Repo) ExistRoleBySlug(slug string) bool {
	exist, err := repo.client.Role.
		Query().
		Where(role.Slug(slug), role.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		repo.log.Error("persist.ExistRoleBySlug()", zap.Error(err))
		return exist
	}

	return exist
}

// CreateRole create role
func (repo *Repo) CreateRole(r *ent.Role, p []*ent.Permission) error {
	// nolint:staticcheck
	r, err := repo.client.Role.
		Create().
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		SetFullAccess(r.FullAccess).
		AddPermissions(p...).
		Save(repo.ctx)

	if err != nil {
		repo.log.Error("persist.CreateRole()", zap.Error(err))
		return err
	}

	return nil
}

// UpdateRole update role
func (repo *Repo) UpdateRole(r *ent.Role, p []*ent.Permission) error {
	_, err := repo.client.Role.
		Update().
		Where(role.ID(r.ID), role.DeleteTimeIsNil()).
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		SetFullAccess(r.FullAccess).
		ClearPermissions().
		AddPermissions(p...).
		Save(repo.ctx)

	if err != nil {
		repo.log.Error("persist.UpdateRole()", zap.Error(err))
		return err
	}

	return nil
}

// SoftDeleteRole soft deletes role
func (repo *Repo) SoftDeleteRole(id uint64) error {
	_, err := repo.client.Role.
		Update().
		Where(role.ID(id), role.DeleteTimeIsNil()).
		SetDeleteTime(time.Now()).
		ClearPermissions().
		Save(repo.ctx)

	if err != nil {
		repo.log.Error("persist.SoftDeleteRole()", zap.Error(err))
		return err
	}

	return nil
}
