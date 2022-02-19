package client

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/google/wire"
	"github.com/spf13/viper"

	// postgresql driver
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/ent/migrate"
	// runtime entgo
	_ "github.com/xdorro/golang-grpc-base-project/api/ent/runtime"
)

// ProviderSet is client providers.
var ProviderSet = wire.NewSet(NewClient)

// NewClient database with config
func NewClient(ctx context.Context, log *zap.Logger) *ent.Client {
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

	// Create an ent.Client.
	client := ent.NewClient(ent.Driver(drv))
	opts := []schema.MigrateOption{
		// migrate.WithGlobalUniqueID(true),
		migrate.WithForeignKeys(false), // Disable foreign keys.
	}

	if viper.GetBool("DEBUG_ENABLE") {
		client = client.Debug()
	}

	if viper.GetBool("DB_MIGRATE") {
		log.Info("Migrating...")

		// Run migration.
		if err = client.Schema.Create(ctx, opts...); err != nil {
			defer func() {
				_ = client.Close()
			}()
			log.Fatal("failed creating schema resources", zap.Error(err))
		}

		log.Info("Migrated")
	}

	return client
}
