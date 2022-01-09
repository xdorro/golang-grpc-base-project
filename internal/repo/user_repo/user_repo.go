package user_repo

import (
	"context"
	"time"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/ent/user"
)

// ProviderSet is repo providers.
var ProviderSet = wire.NewSet(NewRepo)

var _ UserPersist = (*UserRepo)(nil)

type UserRepo struct {
	ctx    context.Context
	client *ent.Client
	log    *zap.Logger
}

func NewRepo(ctx context.Context, client *ent.Client, log *zap.Logger) UserPersist {
	return &UserRepo{
		ctx:    ctx,
		client: client,
		log:    log,
	}
}

// FindAllUsers find all users
func (repo *UserRepo) FindAllUsers() []*ent.User {
	users, err := repo.client.User.
		Query().
		Select(
			user.FieldID,
			user.FieldName,
			user.FieldEmail,
			user.FieldStatus,
		).
		Where(user.DeleteTimeIsNil()).
		All(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindAllUsers()", zap.Error(err))
		return nil
	}

	return users
}

// CreateUser handler CreateUser persist
func (repo *UserRepo) CreateUser(u *ent.User) error {
	u, err := repo.client.User.
		Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		SetRoleID(u.RoleID).
		Save(repo.ctx)

	if err != nil {
		repo.log.Error("persist.CreateUser()", zap.Error(err))
		return err
	}

	return nil
}

// UpdateUser handler UpdateUser persist
func (repo *UserRepo) UpdateUser(u *ent.User) error {
	_, err := repo.client.User.
		Update().
		Where(user.ID(u.ID), user.DeleteTimeIsNil()).
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		SetRoleID(u.RoleID).
		Save(repo.ctx)

	if err != nil {
		repo.log.Error("persist.UpdateUser()", zap.Error(err))
		return err
	}

	return nil
}

func (repo *UserRepo) FindUserByEmail(email string) (*ent.User, error) {
	obj, err := repo.client.User.
		Query().
		Where(user.Email(email), user.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindUserByEmail()", zap.Error(err))
		return nil, err
	}

	return obj, nil
}

// FindUserByID handler FindUserByID persist
func (repo *UserRepo) FindUserByID(id uint64) (*ent.User, error) {
	obj, err := repo.client.User.
		Query().
		Where(user.ID(id), user.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		repo.log.Error("persist.FindUserByID()", zap.Error(err))
		return nil, err
	}

	return obj, nil
}

// ExistUserByID return true if ID existed
func (repo *UserRepo) ExistUserByID(id uint64) bool {
	check, err := repo.client.User.
		Query().
		Where(user.ID(id), user.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		repo.log.Error("persist.ExistUserByID()", zap.Error(err))
		return check
	}

	return check
}

// ExistUserByEmail return true if email existed
func (repo *UserRepo) ExistUserByEmail(email string) bool {
	check, err := repo.client.User.
		Query().
		Where(user.Email(email), user.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		repo.log.Error("persist.ExistUserByEmail()", zap.Error(err))
		return check
	}

	return check
}

// DeleteUser delete user by ID
func (repo *UserRepo) DeleteUser(id uint64) error {
	if _, err := repo.client.User.
		Delete().
		Where(user.ID(id)).
		Exec(repo.ctx); err != nil {
		repo.log.Error("persist.DeleteUser()", zap.Error(err))
		return err
	}

	return nil
}

// SoftDeleteUser update user delete time by ID
func (repo *UserRepo) SoftDeleteUser(id uint64) error {
	if _, err := repo.client.User.
		Update().
		Where(user.ID(id), user.DeleteTimeIsNil()).
		SetDeleteTime(time.Now()).
		Save(repo.ctx); err != nil {
		repo.log.Error("persist.SoftDeleteUser()", zap.Error(err))
		return err
	}

	return nil
}
