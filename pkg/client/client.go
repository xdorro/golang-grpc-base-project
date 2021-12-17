package client

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/spf13/viper"

	// postgresql driver
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/ent/migrate"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"

	// runtime entgo
	_ "github.com/xdorro/golang-grpc-base-project/ent/runtime"
	"github.com/xdorro/golang-grpc-base-project/internal/persist"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

type Client struct {
	DB      *ent.Client
	Persist persist.Persist
}

// NewClient database with config
func NewClient(ctx context.Context) *Client {
	driver := viper.GetString("DB_DRIVER")
	url := viper.GetString("DB_URL")

	logger.Info("Connect to database",
		zap.String("driver", driver),
		zap.String("url", url),
	)

	// Open the database connection.
	drv, err := sql.Open(driver, url)
	if err != nil {
		logger.Fatal("sql.Open()", zap.Error(err))
	}

	// Create an ent.client.
	db := ent.NewClient(ent.Driver(drv))

	if viper.GetBool("DB_MIGRATE") {
		logger.Info("Migrating...")
		// Run migration.
		if err = db.Schema.Create(
			ctx,
			migrate.WithGlobalUniqueID(true),
			migrate.WithForeignKeys(false), // Disable foreign keys.
		); err != nil {
			_ = db.Close()
			logger.Fatal("failed creating schema resources", zap.Error(err))
		}

		logger.Info("Migrated")
	}

	// Create new persist
	per := repo.NewRepo(ctx, db)

	client := &Client{
		DB:      db,
		Persist: per,
	}

	return client
}

func (client *Client) Close() error {
	return client.DB.Close()
}
