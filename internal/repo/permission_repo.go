package repo

import (
	"time"

	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/ent/permission"
	"github.com/xdorro/golang-grpc-base-project/ent/role"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

// FindAllPermissions find all permissions
func (repo *Repo) FindAllPermissions() []*ent.Permission {
	permissions, err := repo.client.Permission.
		Query().
		Select(
			permission.FieldID,
			permission.FieldName,
			permission.FieldSlug,
			permission.FieldStatus,
		).
		Where(permission.DeleteTimeIsNil()).
		All(repo.ctx)

	if err != nil {
		logger.Error("persist.FindAllPermissions()", zap.Error(err))
		return nil
	}

	return permissions
}

// FindPermissionByID find permission by ID
func (repo *Repo) FindPermissionByID(id uint64) (*ent.Permission, error) {
	r, err := repo.client.Permission.
		Query().
		Where(permission.ID(id), permission.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		logger.Error("persist.FindPermissionByID()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindPermissionBySlug find permission by slug
func (repo *Repo) FindPermissionBySlug(slug string) (*ent.Permission, error) {
	r, err := repo.client.Permission.
		Query().
		Where(permission.SlugEqualFold(slug), permission.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		logger.Error("persist.FindPermissionBySlug()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindPermissionByIDAndRoleIDNot find permission by ID and roleID not
func (repo *Repo) FindPermissionByIDAndRoleIDNot(id uint64, roleId uint64) (*ent.Permission, error) {
	p, err := repo.client.Permission.
		Query().
		Where(
			permission.ID(id),
			permission.DeleteTimeIsNil(),
			permission.HasRolesWith(role.Not(role.ID(roleId))),
		).
		First(repo.ctx)

	if err != nil {
		logger.Error("persist.FindPermissionByIDAndRoleIDNot()", zap.Error(err))
		return nil, err
	}

	return p, nil
}

// ExistPermissionByID exists a permission by ID
func (repo *Repo) ExistPermissionByID(id uint64) bool {
	exist, err := repo.client.Permission.
		Query().
		Where(permission.ID(id), permission.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		logger.Error("persist.ExistPermissionByID()", zap.Error(err))
		return exist
	}

	return exist
}

// ExistPermissionBySlug exist permission by slug
func (repo *Repo) ExistPermissionBySlug(slug string) bool {
	exist, err := repo.client.Permission.
		Query().
		Where(permission.Slug(slug), permission.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		logger.Error("persist.ExistPermissionBySlug()", zap.Error(err))
		return exist
	}

	return exist
}

// CreatePermission create permission
func (repo *Repo) CreatePermission(r *ent.Permission) error {
	r, err := repo.client.Permission.
		Create().
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		Save(repo.ctx)

	if err != nil {
		logger.Error("persist.CreatePermission()", zap.Error(err))
		return err
	}

	return nil
}

// CreatePermissionBulk create permission bulk
func (repo *Repo) CreatePermissionBulk(r []*ent.PermissionCreate) error {
	_, err := repo.client.Permission.
		CreateBulk(r...).
		Save(repo.ctx)

	if err != nil {
		logger.Error("persist.CreatePermissionBulk()", zap.Error(err))
		return err
	}

	return nil
}

// UpdatePermission update permission
func (repo *Repo) UpdatePermission(r *ent.Permission) error {
	_, err := repo.client.Permission.
		Update().
		Where(permission.ID(r.ID), permission.DeleteTimeIsNil()).
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		Save(repo.ctx)

	if err != nil {
		logger.Error("persist.UpdatePermission()", zap.Error(err))
		return err
	}

	return nil
}

// SoftDeletePermission soft delete permission
func (repo *Repo) SoftDeletePermission(id uint64) error {
	if _, err := repo.client.Permission.
		Update().
		Where(permission.ID(id), permission.DeleteTimeIsNil()).
		SetDeleteTime(time.Now()).
		ClearRoles().
		Save(repo.ctx); err != nil {
		logger.Error("persist.SoftDeletePermission()", zap.Error(err))
		return err
	}

	return nil
}
