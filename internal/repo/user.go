package repo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/internal/models"
)

// FindAllUsers returns all users
func (r *Repo) FindAllUsers(filter any, opt ...*options.FindOptions) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 1*time.Minute)
	defer cancel()

	cur, err := r.
		userCollection().
		Find(ctx, filter, opt...)
	if err != nil {
		r.log.Error("Error find all users", zap.Error(err))
		return nil, err
	}

	defer func() {
		_ = cur.Close(ctx)
	}()

	var data []*models.User
	for cur.Next(ctx) {
		user := &models.User{}
		if err = cur.Decode(user); err != nil {
			r.log.Error("Error find all users", zap.Error(err))
			return nil, err
		}

		data = append(data, user)
	}

	return data, nil
}

// CountAllUsers returns total number of users
func (r *Repo) CountAllUsers(filter any) (int64, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 1*time.Minute)
	defer cancel()

	total, err := r.
		userCollection().
		CountDocuments(ctx, filter)
	if err != nil {
		r.log.Error("Error count all users", zap.Error(err))
		return 0, err
	}

	return total, nil
}

// FindUser returns user
func (r *Repo) FindUser(filter any, opt ...*options.FindOneOptions) (*models.User, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 1*time.Minute)
	defer cancel()

	result := &models.User{}
	err := r.
		userCollection().
		FindOne(ctx, filter, opt...).
		Decode(result)
	if err != nil {
		r.log.Error("Error find user", zap.Error(err))
		return nil, err
	}

	return result, nil
}

// CreateUser creates user
func (r *Repo) CreateUser(data any, opt ...*options.InsertOneOptions) error {
	ctx, cancel := context.WithTimeout(r.ctx, 1*time.Minute)
	defer cancel()

	_, err := r.
		userCollection().
		InsertOne(ctx, data, opt...)
	if err != nil {
		r.log.Error("Error creating user", zap.Error(err))
		return err
	}

	return nil
}

// UpdateUser updates user
func (r *Repo) UpdateUser(filter, data any, opt ...*options.UpdateOptions) error {
	ctx, cancel := context.WithTimeout(r.ctx, 1*time.Minute)
	defer cancel()

	_, err := r.
		userCollection().
		UpdateOne(ctx, filter, data, opt...)
	if err != nil {
		r.log.Error("Error updating user", zap.Error(err))
		return err
	}

	return nil
}

// DeleteUser deletes user
func (r *Repo) DeleteUser(filter any, opt ...*options.DeleteOptions) error {
	ctx, cancel := context.WithTimeout(r.ctx, 1*time.Minute)
	defer cancel()

	_, err := r.
		userCollection().
		DeleteOne(ctx, filter, opt...)
	if err != nil {
		r.log.Error("Error deleting user", zap.Error(err))
		return err
	}

	return nil
}
