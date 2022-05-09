package config

import (
	"github.com/spf13/viper"
)

// defaultConfig is the default configuration for the application.
func defaultConfig() {
	// APP
	viper.SetDefault("APP_NAME", "Golang gRPC Base Project")
	viper.SetDefault("APP_VERSION", "0.0.0")
	viper.SetDefault("APP_PORT", "8088")
	viper.SetDefault("APP_KEY", "./config/cert/key.pem")
	viper.SetDefault("APP_CERT", "./config/cert/cert.pem")
	viper.SetDefault("APP_SECRET_KEY", "your-256-bit-secret")
	viper.SetDefault("LOG_PAYLOAD", true)

	// DATABASE
	viper.SetDefault("DB_URL", "mongodb://localhost:27017")
	viper.SetDefault("DB_NAME", "base")

	// REDIS
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
}
