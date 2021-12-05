package repo

import (
	"time"

	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base/pkg/ent"
	"github.com/kucow/golang-grpc-base/pkg/ent/user"
)

func (repo *Repo) FindAllUsers() []*ent.User {
	users, err := repo.Client.User.
		Query().
		Select(
			user.FieldID,
			user.FieldName,
			user.FieldEmail,
			user.FieldStatus,
		).
		Where(user.DeleteTimeIsNil()).
		All(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindAllUsers()", zap.Error(err))
		return nil
	}

	return users
}

// CreateUser handler CreateUser persist
func (repo *Repo) CreateUser(u *ent.User, r []*ent.Role) error {
	u, err := repo.Client.User.
		Create().
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		AddRoles(r...).
		Save(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.CreateUser()", zap.Error(err))
		return err
	}

	return nil
}

// UpdateUser handler UpdateUser persist
func (repo *Repo) UpdateUser(u *ent.User, r []*ent.Role) error {
	_, err := repo.Client.User.
		Update().
		Where(user.ID(u.ID), user.DeleteTimeIsNil()).
		SetName(u.Name).
		SetEmail(u.Email).
		SetPassword(u.Password).
		SetStatus(u.Status).
		ClearRoles().
		AddRoles(r...).
		Save(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.UpdateUser()", zap.Error(err))
		return err
	}

	return nil
}

// FindUserByEmail handler FindUserByEmail persist
func (repo *Repo) FindUserByEmail(email string) (*ent.User, error) {
	obj, err := repo.Client.User.
		Query().
		Where(user.Email(email), user.DeleteTimeIsNil()).
		First(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindUserByEmail()", zap.Error(err))
		return nil, err
	}

	return obj, nil
}

// FindUserByID handler FindUserByID persist
func (repo *Repo) FindUserByID(id string) (*ent.User, error) {
	obj, err := repo.Client.User.
		Query().
		Where(user.ID(id), user.DeleteTimeIsNil()).
		First(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.FindUserByID()", zap.Error(err))
		return nil, err
	}

	return obj, nil
}

// ExistUserByEmail return true if email existed
func (repo *Repo) ExistUserByEmail(email string) bool {
	check, err := repo.Client.User.
		Query().
		Where(user.Email(email), user.DeleteTimeIsNil()).
		Exist(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.ExistUserByEmail()", zap.Error(err))
		return check
	}

	return check
}

// ExistUserByID return true if ID existed
func (repo *Repo) ExistUserByID(id string) bool {
	check, err := repo.Client.User.
		Query().
		Where(user.ID(id), user.DeleteTimeIsNil()).
		Exist(repo.Ctx)

	if err != nil {
		repo.Log.Error("persist.ExistUserByID()", zap.Error(err))
		return check
	}

	return check
}

// DeleteUser delete user by ID
func (repo *Repo) DeleteUser(id string) error {
	if _, err := repo.Client.User.
		Delete().
		Where(user.ID(id)).
		Exec(repo.Ctx); err != nil {
		repo.Log.Error("persist.DeleteUser()", zap.Error(err))
		return err
	}

	return nil
}

// SoftDeleteUser update user delete time by ID
func (repo *Repo) SoftDeleteUser(id string) error {
	if _, err := repo.Client.User.
		Update().
		Where(user.ID(id), user.DeleteTimeIsNil()).
		SetDeleteTime(time.Now()).
		Save(repo.Ctx); err != nil {
		repo.Log.Error("persist.SoftDeleteUser()", zap.Error(err))
		return err
	}

	return nil
}
