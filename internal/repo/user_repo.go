package repo

import (
	"time"

	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/ent/user"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

// FindAllUsers find all users
func (repo *Repo) FindAllUsers() []*ent.User {
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
		logger.Error("persist.FindAllUsers()", zap.Error(err))
		return nil
	}

	return users
}

// CreateUser handler CreateUser persist
func (repo *Repo) CreateUser(u *ent.User, r []*ent.Role) error {
	u, err := repo.client.User.
		Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		AddRoles(r...).
		Save(repo.ctx)

	if err != nil {
		logger.Error("persist.CreateUser()", zap.Error(err))
		return err
	}

	return nil
}

// UpdateUser handler UpdateUser persist
func (repo *Repo) UpdateUser(u *ent.User, r []*ent.Role) error {
	_, err := repo.client.User.
		Update().
		Where(user.ID(u.ID), user.DeleteTimeIsNil()).
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		ClearRoles().
		AddRoles(r...).
		Save(repo.ctx)

	if err != nil {
		logger.Error("persist.UpdateUser()", zap.Error(err))
		return err
	}

	return nil
}

// FindUserByEmail handler FindUserByEmail persist
func (repo *Repo) FindUserByEmail(email string) (*ent.User, error) {
	obj, err := repo.client.User.
		Query().
		Where(user.Email(email), user.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		logger.Error("persist.FindUserByEmail()", zap.Error(err))
		return nil, err
	}

	return obj, nil
}

// FindUserByID handler FindUserByID persist
func (repo *Repo) FindUserByID(id uint64) (*ent.User, error) {
	obj, err := repo.client.User.
		Query().
		Where(user.ID(id), user.DeleteTimeIsNil()).
		First(repo.ctx)

	if err != nil {
		logger.Error("persist.FindUserByID()", zap.Error(err))
		return nil, err
	}

	return obj, nil
}

// ExistUserByEmail return true if email existed
func (repo *Repo) ExistUserByEmail(email string) bool {
	check, err := repo.client.User.
		Query().
		Where(user.Email(email), user.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		logger.Error("persist.ExistUserByEmail()", zap.Error(err))
		return check
	}

	return check
}

// ExistUserByID return true if ID existed
func (repo *Repo) ExistUserByID(id uint64) bool {
	check, err := repo.client.User.
		Query().
		Where(user.ID(id), user.DeleteTimeIsNil()).
		Exist(repo.ctx)

	if err != nil {
		logger.Error("persist.ExistUserByID()", zap.Error(err))
		return check
	}

	return check
}

// DeleteUser delete user by ID
func (repo *Repo) DeleteUser(id uint64) error {
	if _, err := repo.client.User.
		Delete().
		Where(user.ID(id)).
		Exec(repo.ctx); err != nil {
		logger.Error("persist.DeleteUser()", zap.Error(err))
		return err
	}

	return nil
}

// SoftDeleteUser update user delete time by ID
func (repo *Repo) SoftDeleteUser(id uint64) error {
	if _, err := repo.client.User.
		Update().
		Where(user.ID(id), user.DeleteTimeIsNil()).
		SetDeleteTime(time.Now()).
		Save(repo.ctx); err != nil {
		logger.Error("persist.SoftDeleteUser()", zap.Error(err))
		return err
	}

	return nil
}
