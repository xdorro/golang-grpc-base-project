package role_repo

import (
	"context"
	"strings"
	"time"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/ent/permission"
	"github.com/xdorro/golang-grpc-base-project/api/ent/role"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(NewRepo)

var _ RolePersist = (*RoleRepo)(nil)

type RoleRepo struct {
	ctx    context.Context
	client *ent.Client
	log    *zap.Logger
}

func NewRepo(ctx context.Context, client *ent.Client, log *zap.Logger) RolePersist {
	return &RoleRepo{
		ctx:    ctx,
		client: client,
		log:    log,
	}
}

// FindAllRoles find all roles
func (repo *RoleRepo) FindAllRoles() []*ent.Role {
	roles, err := repo.client.Role.
		Query().
		Select(
			role.FieldID,
			role.FieldName,
			role.FieldSlug,
			role.FieldFullAccess,
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
func (repo *RoleRepo) FindRoleByID(id uint64) (*ent.Role, error) {
	r, err := repo.client.Role.
		Query().
		Select(
			role.FieldID,
			role.FieldName,
			role.FieldSlug,
			role.FieldFullAccess,
			role.FieldStatus,
		).
		Where(role.ID(id), role.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindRoleByID()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindRoleBySlug find role by slug
func (repo *RoleRepo) FindRoleBySlug(slug string) (*ent.Role, error) {
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
func (repo *RoleRepo) FindRoleByIDAndPermissionID(id, permissionID uint64) (*ent.Role, error) {
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
func (repo *RoleRepo) FindRoleByIDAndPermissionIDNot(id, permissionID uint64) (*ent.Role, error) {
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
func (repo *RoleRepo) ExistRoleByID(id uint64) bool {
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
func (repo *RoleRepo) ExistRoleBySlug(slug string) bool {
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

// ExistRoleByIDNotSlug exist role by slug
func (repo *RoleRepo) ExistRoleByIDNotSlug(id uint64, slug string) bool {
	exist, err := repo.client.Role.
		Query().
		Where(role.IDNEQ(id), role.Slug(slug), role.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		repo.log.Error("persist.ExistRoleByIDNotSlug()", zap.Error(err))
		return exist
	}

	return exist
}

// CreateRole create role
func (repo *RoleRepo) CreateRole(r *ent.Role, p []*ent.Permission) error {
	// nolint:staticcheck
	r, err := repo.client.Role.
		Create().
		SetName(strings.TrimSpace(r.Name)).
		SetSlug(strings.TrimSpace(r.Slug)).
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
func (repo *RoleRepo) UpdateRole(r *ent.Role, p []*ent.Permission) error {
	_, err := repo.client.Role.
		Update().
		Where(role.ID(r.ID), role.DeleteTimeIsNil()).
		SetName(strings.TrimSpace(r.Name)).
		SetSlug(strings.TrimSpace(r.Slug)).
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
func (repo *RoleRepo) SoftDeleteRole(id uint64) error {
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
