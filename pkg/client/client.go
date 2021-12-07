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
func NewClient(ctx context.Context, log *zap.Logger) *Client {
	driver := viper.GetString("DB_DRIVER")
	url := viper.GetString("DB_URL")

	log.Info("Connect to database",
		zap.String("driver", driver),
		zap.String("url", url),
	)

	// Open the database connection.
	drv, err := sql.Open(driver, url)
	if err != nil {
		log.Fatal("sql.Open()", zap.Error(err))
	}

	// Create an ent.client.
	db := ent.NewClient(ent.Driver(drv))

	if viper.GetBool("DB_MIGRATE") {
		log.Info("Migrating...")
		// Run migration.
		if err = db.Schema.Create(
			ctx,
			migrate.WithGlobalUniqueID(true),
		); err != nil {
			_ = db.Close()
			log.Fatal("failed creating schema resources", zap.Error(err))
		}

		log.Info("Migrated")
	}

	client := &Client{
		DB: db,
		// Create new persist
		Persist: repo.NewRepo(ctx, log, db),
	}

	return client
}

func (client *Client) Close() error {
	return client.DB.Close()
}
