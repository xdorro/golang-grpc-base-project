package permission_repo

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

var _ PermissionPersist = (*PermissionRepo)(nil)

type PermissionRepo struct {
	ctx    context.Context
	client *ent.Client
	log    *zap.Logger
}

func NewRepo(ctx context.Context, client *ent.Client, log *zap.Logger) PermissionPersist {
	return &PermissionRepo{
		ctx:    ctx,
		client: client,
		log:    log,
	}
}

// FindAllPermissions find all permissions
func (repo *PermissionRepo) FindAllPermissions() []*ent.Permission {
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
		repo.log.Error("persist.FindAllPermissions()", zap.Error(err))
		return nil
	}

	return permissions
}

// FindAllPermissionsWithRoles find all permissions
func (repo *PermissionRepo) FindAllPermissionsWithRoles() []*ent.Permission {
	permissions, err := repo.client.Permission.
		Query().
		Select(
			permission.FieldID,
			permission.FieldName,
			permission.FieldSlug,
			permission.FieldStatus,
		).
		Where(permission.DeleteTimeIsNil()).
		WithRoles(func(q *ent.RoleQuery) {
			q.Where(role.DeleteTimeIsNil())
		}).
		All(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindAllPermissions()", zap.Error(err))
		return nil
	}

	return permissions
}

func (repo *PermissionRepo) FindAllPermissionBySlugs(slugs []string) []*ent.Permission {
	permissions, err := repo.client.Permission.
		Query().
		Select(
			permission.FieldID,
			permission.FieldName,
			permission.FieldSlug,
			permission.FieldStatus,
		).
		Where(permission.SlugIn(slugs...), permission.DeleteTimeIsNil()).
		All(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindAllPermissionBySlugs()", zap.Error(err))
		return nil
	}

	return permissions
}

// FindPermissionByID find permission by ID
func (repo *PermissionRepo) FindPermissionByID(id uint64) (*ent.Permission, error) {
	r, err := repo.client.Permission.
		Query().
		Where(permission.ID(id), permission.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindPermissionByID()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindPermissionBySlug find permission by slug
func (repo *PermissionRepo) FindPermissionBySlug(slug string) (*ent.Permission, error) {
	r, err := repo.client.Permission.
		Query().
		Where(permission.SlugEqualFold(slug), permission.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindPermissionBySlug()", zap.Error(err))
		return nil, err
	}

	return r, nil
}

// FindPermissionByIDAndRoleIDNot find permission by ID and roleID not
func (repo *PermissionRepo) FindPermissionByIDAndRoleIDNot(id uint64, roleId uint64) (*ent.Permission, error) {
	p, err := repo.client.Permission.
		Query().
		Where(
			permission.ID(id),
			permission.DeleteTimeIsNil(),
			permission.HasRolesWith(role.Not(role.ID(roleId))),
		).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindPermissionByIDAndRoleIDNot()", zap.Error(err))
		return nil, err
	}

	return p, nil
}

// ExistPermissionByID exists a permission by ID
func (repo *PermissionRepo) ExistPermissionByID(id uint64) bool {
	exist, err := repo.client.Permission.
		Query().
		Where(permission.ID(id), permission.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		repo.log.Error("persist.ExistPermissionByID()", zap.Error(err))
		return exist
	}

	return exist
}

// ExistPermissionBySlug exist permission by slug
func (repo *PermissionRepo) ExistPermissionBySlug(slug string) bool {
	exist, err := repo.client.Permission.
		Query().
		Where(permission.Slug(slug), permission.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		repo.log.Error("persist.ExistPermissionBySlug()", zap.Error(err))
		return exist
	}

	return exist
}

// ExistPermissionByIDNotAndSlug exist permission by slug
func (repo *PermissionRepo) ExistPermissionByIDNotAndSlug(id uint64, slug string) bool {
	exist, err := repo.client.Permission.
		Query().
		Where(permission.IDNEQ(id), permission.Slug(slug), permission.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		repo.log.Error("persist.ExistPermissionByIDNotAndSlug()", zap.Error(err))
		return exist
	}

	return exist
}

// CreatePermission create permission
func (repo *PermissionRepo) CreatePermission(r *ent.Permission) error {
	r, err := repo.client.Permission.
		Create().
		SetName(strings.TrimSpace(r.Name)).
		SetSlug(strings.TrimSpace(r.Slug)).
		SetStatus(r.Status).
		Save(repo.ctx)

	if err != nil {
		repo.log.Error("persist.CreatePermission()", zap.Error(err))
		return err
	}

	return nil
}

// CreatePermissionBulk create permission bulk
func (repo *PermissionRepo) CreatePermissionBulk(r []*ent.PermissionCreate) error {
	_, err := repo.client.Permission.
		CreateBulk(r...).
		Save(repo.ctx)

	if err != nil {
		repo.log.Error("persist.CreatePermissionBulk()", zap.Error(err))
		return err
	}

	return nil
}

// UpdatePermission update permission
func (repo *PermissionRepo) UpdatePermission(r *ent.Permission) error {
	_, err := repo.client.Permission.
		Update().
		Where(permission.ID(r.ID), permission.DeleteTimeIsNil()).
		SetName(strings.TrimSpace(r.Name)).
		SetSlug(strings.TrimSpace(r.Slug)).
		SetStatus(r.Status).
		Save(repo.ctx)

	if err != nil {
		repo.log.Error("persist.UpdatePermission()", zap.Error(err))
		return err
	}

	return nil
}

// SoftDeletePermission soft delete permission
func (repo *PermissionRepo) SoftDeletePermission(id uint64) error {
	if _, err := repo.client.Permission.
		Update().
		Where(permission.ID(id), permission.DeleteTimeIsNil()).
		SetDeleteTime(time.Now()).
		ClearRoles().
		Save(repo.ctx); err != nil {
		repo.log.Error("persist.SoftDeletePermission()", zap.Error(err))
		return err
	}

	return nil
}
