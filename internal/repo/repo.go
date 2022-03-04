package repo

import (
	"context"
	"time"

	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/internal/models"
)

// ProviderRepoSet is repository providers.
var ProviderRepoSet = wire.NewSet(NewRepo)
var _ IRepo = (*Repo)(nil)

// IRepo is the interface for all repositories
type IRepo interface {
	Close() error
	Collection(collectionName string) *mongo.Collection
	UserCollection() *mongo.Collection
}

// Repo is repository struct.
type Repo struct {
	ctx    context.Context
	log    *zap.Logger
	client *mongo.Client
}

// NewRepo creates new repository.
func NewRepo(ctx context.Context, log *zap.Logger) IRepo {
	uri := viper.GetString("MONGODB_URI")

	log.Info("Connecting to MongoDB", zap.String("uri", uri))

	clientOpts := options.Client().
		ApplyURI(uri)

	// Set connect timeout to 15 seconds
	ctxConn, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Create new client and connect to MongoDB
	client, err := mongo.Connect(ctxConn, clientOpts)
	if err != nil {
		log.Panic("Failed to connect to mongodb", zap.Error(err))
	}

	// Set the ping timeout to 5 seconds
	ctxPing, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	// Ping the primary
	if err = client.Ping(ctxPing, readpref.Primary()); err != nil {
		log.Panic("Failed to Ping() to mongodb", zap.Error(err))
	}

	log.Info("Successfully connected and pinged.")

	return &Repo{
		ctx:    ctx,
		log:    log,
		client: client,
	}
}

// Close closes the repository.
func (r *Repo) Close() error {
	if err := r.client.Disconnect(r.ctx); err != nil {
		r.log.Error("Failed to disconnect from mongodb", zap.Error(err))
		return err
	}

	return nil
}

// Collection returns the mongo collection by Name.
func (r *Repo) Collection(collectionName string) *mongo.Collection {
	dbName := viper.GetString("MONGODB_DBNAME")

	return r.client.Database(dbName).Collection(collectionName)
}

// UserCollection returns the mongo collection for users.
func (r *Repo) UserCollection() *mongo.Collection {
	return r.Collection(models.User{}.CollectionName())
}
