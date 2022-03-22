package repo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/xdorro/golang-grpc-base-project/internal/models"
)

// IRepo is the interface for all repositories
type IRepo interface {
	Close() error
	Collection(collectionName string) *mongo.Collection

	IUserRepo
}

// IUserRepo is the interface for user repositories
type IUserRepo interface {
	FindAllUsers(filter any, opt ...*options.FindOptions) ([]*models.User, error)
	CountAllUsers(filter any) (int64, error)
	FindUser(filter any, opt ...*options.FindOneOptions) (*models.User, error)
	CreateUser(data any, opt ...*options.InsertOneOptions) error
	UpdateUser(filter, data any, opt ...*options.UpdateOptions) error
	DeleteUser(filter any, opt ...*options.DeleteOptions) error
}
