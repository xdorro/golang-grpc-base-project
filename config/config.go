package config

import (
	"errors"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

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

	// Set default values
	defaultConfig()

	// Replace env key
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	log.Info(viper.GetString("APP_NAME"),
		zap.String("app-version", viper.GetString("APP_VERSION")),
		zap.String("go-version", runtime.Version()),
	)
}
