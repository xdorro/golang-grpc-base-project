package repo

import (
	"time"

	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base/pkg/ent"
	"github.com/kucow/golang-grpc-base/pkg/ent/permission"
	"github.com/kucow/golang-grpc-base/pkg/ent/role"
)

func (repo *Repo) FindAllPermissions() []*ent.Permission {
	permissions, err := repo.Client.Permission.
		Query().
		Select(
			permission.FieldID,
			permission.FieldName,
			permission.FieldSlug,
			permission.FieldStatus,
		).
		Where(permission.DeleteTimeIsNil()).
		All(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindAllPermissions()", zap.Error(err))
		return nil
	}

	return permissions
}

func (repo *Repo) FindPermissionByID(id string) (*ent.Permission, error) {
	r, err := repo.Client.Permission.
		Query().
		Where(permission.ID(id), permission.DeleteTimeIsNil()).
		First(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindPermissionByID()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

func (repo *Repo) FindPermissionByIDAndRoleIDNot(id string, roleId string) (*ent.Permission, error) {
	p, err := repo.Client.Permission.
		Query().
		Where(
			permission.ID(id),
			permission.DeleteTimeIsNil(),
			permission.HasRolesWith(role.Not(role.ID(roleId))),
		).
		First(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindPermissionByIDAndRoleIDNot()", zap.Error(err))
		return nil, err
	}

	return p, nil
}

func (repo *Repo) ExistPermissionByID(id string) bool {
	exist, err := repo.Client.Permission.
		Query().
		Where(permission.ID(id), permission.DeleteTimeIsNil()).
		Exist(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.ExistPermissionByID()", zap.Error(err))
		return exist
	}

	return exist
}

func (repo *Repo) ExistPermissionBySlug(slug string) bool {
	exist, err := repo.Client.Permission.
		Query().
		Where(permission.Slug(slug), permission.DeleteTimeIsNil()).
		Exist(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.ExistPermissionBySlug()", zap.Error(err))
		return exist
	}

	return exist
}

func (repo *Repo) CreatePermission(r *ent.Permission) error {
	r, err := repo.Client.Permission.
		Create().
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		Save(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.CreatePermission()", zap.Error(err))
		return err
	}

	return nil
}

func (repo *Repo) UpdatePermission(r *ent.Permission) error {
	_, err := repo.Client.Permission.
		Update().
		Where(permission.ID(r.ID), permission.DeleteTimeIsNil()).
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		Save(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.UpdatePermission()", zap.Error(err))
		return err
	}

	return nil
}

func (repo *Repo) SoftDeletePermission(id string) error {
	if _, err := repo.Client.Permission.
		Update().
		Where(permission.ID(id), permission.DeleteTimeIsNil()).
		SetDeleteTime(time.Now()).
		ClearRoles().
		Save(repo.Ctx); err != nil {
		repo.Log.Error("persist.SoftDeletePermission()", zap.Error(err))
		return err
	}

	return nil
}
