package client

import (
	"entgo.io/ent/dialect/sql"
	"github.com/spf13/viper"

	// postgresql driver
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base-project/internal/common"
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	_ "github.com/kucow/golang-grpc-base-project/pkg/ent/runtime"
)

// NewClient database with config
func NewClient(opts *common.Option) {
	driver := viper.GetString("DB_DRIVER")
	url := viper.GetString("DB_URL")

	opts.Log.Info("Connect to database",
		zap.String("driver", driver),
		zap.String("url", url),
	)

	// Open the database connection.
	db, err := sql.Open(driver, url)
	if err != nil {
		opts.Log.Fatal("sql.Open()", zap.Error(err))
	}

	// Create an ent.Client.
	opts.Client = ent.NewClient(ent.Driver(db))

	if viper.GetBool("DB_MIGRATE") {
		opts.Log.Info("Migrating...")
		// Run migration.
		if err = opts.Client.Schema.Create(
			opts.Ctx,
			// migrate.WithGlobalUniqueID(true),
		); err != nil {
			_ = opts.Client.Close()
			opts.Log.Fatal("failed creating schema resources", zap.Error(err))
		}

		opts.Log.Info("Migrated")
	}
}
