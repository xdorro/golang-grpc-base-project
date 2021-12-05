package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

// loadDefault load default config
func loadDefault() {
	viper.SetDefault("APP_NAME", "Golang gRPC Base Project")
	viper.SetDefault("APP_VERSION", "0.0.0")
	viper.SetDefault("GRPC_PORT", "3100")
	viper.SetDefault("REST_PORT", "3200")
	viper.SetDefault("LOG_PAYLOAD", true)
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
func NewConfig(path ...string) {
	configPath := "."
	if len(path) > 0 {
		configPath = path[0]
	}

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.AddConfigPath(configPath)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found; ignore error if desired
		log.Printf("viper.ReadInConfig(): %v", err)
	}

	viper.AutomaticEnv()

	// Load defaultConfig
	loadDefault()

	// Replace env key
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
