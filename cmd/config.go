package main

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

// loadDefault load default config
func loadDefault() {
	// APP
	viper.SetDefault("APP_NAME", "Golang gRPC Base Project")
	viper.SetDefault("APP_VERSION", "0.0.0")
	viper.SetDefault("APP_DEBUG", false)
	viper.SetDefault("GRPC_PORT", "3100")
	viper.SetDefault("HTTP_PORT", "3200")

	viper.SetDefault("LOG_PAYLOAD", true)
	viper.SetDefault("SEEDER_SERVICE", false)
	viper.SetDefault("METRIC_ENABLE", false)
	// ASYNQ
	viper.SetDefault("ASYNQ_ENABLE", false)

	viper.SetDefault("MACHINE_ID", "0")
	viper.SetDefault("AUTH_SECRET_KEY", "your-256-bit-secret")
	// DATABASE
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_MIGRATE", true)
	// REDIS
	viper.SetDefault("REDIS_URL", "localhost:16379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
}

// NewConfig create new config
func init() {
	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found; ignore error if desired
		logger.Error("viper.ReadInConfig()", zap.Error(err))
	}

	viper.AutomaticEnv()

	// Load defaultConfig
	loadDefault()

	// Replace env key
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
