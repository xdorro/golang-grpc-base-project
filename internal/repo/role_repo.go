package repo

import (
	"time"

	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	"github.com/kucow/golang-grpc-base-project/pkg/ent/permission"
	"github.com/kucow/golang-grpc-base-project/pkg/ent/role"
)

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

func (repo *Repo) CreateRole(r *ent.Role, p []*ent.Permission) error {
	r, err := repo.Client.Role.
		Create().
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		AddPermissions(p...).
		Save(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.CreateRole()", zap.Error(err))
		return err
	}

	return nil
}

func (repo *Repo) UpdateRole(r *ent.Role, p []*ent.Permission) error {
	_, err := repo.Client.Role.
		Update().
		Where(role.ID(r.ID), role.DeleteTimeIsNil()).
		SetName(r.Name).
		SetSlug(r.Slug).
		SetStatus(r.Status).
		ClearPermissions().
		AddPermissions(p...).
		Save(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.UpdateRole()", zap.Error(err))
		return err
	}

	return nil
}

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
