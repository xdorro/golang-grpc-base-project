package config

import (
	"errors"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// loadDefault load default config
func loadDefault() {
	// APP
	viper.SetDefault("APP_NAME", "Golang gRPC Base Project")
	viper.SetDefault("APP_VERSION", "0.0.0")
	viper.SetDefault("GRPC_PORT", "3100")
	viper.SetDefault("HTTP_PORT", "3200")
	// MODE ENABLE
	viper.SetDefault("DEBUG_ENABLE", false)
	viper.SetDefault("LOG_PAYLOAD", true)
	viper.SetDefault("SEEDER_SERVICE", false)

	viper.SetDefault("MACHINE_ID", "0")
	viper.SetDefault("AUTH_SECRET_KEY", "your-256-bit-secret")
	// DATABASE
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_MIGRATE", true)
	// REDIS
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
}

// NewConfig create new config
func NewConfig(log *zap.Logger) {
	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found; ignore error if desired
		var notfound viper.ConfigFileNotFoundError
		if ok := errors.Is(err, notfound); !ok {
			log.Error("viper.ReadInConfig()", zap.Error(err))
		}
	}

	viper.AutomaticEnv()

	// Load defaultConfig
	loadDefault()

	// Replace env key
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	log.Info(viper.GetString("APP_NAME"),
		zap.String("app-version", viper.GetString("APP_VERSION")),
		zap.String("go-version", runtime.Version()),
	)
}
