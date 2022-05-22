package config

import (
	"errors"
	"runtime"
	"strings"

	"github.com/spf13/viper"

	"github.com/xdorro/golang-grpc-base-project/internal/log"
)

// InitConfig create new config
func InitConfig() {
	log.Info().
		Str("goarch", runtime.GOARCH).
		Str("goos", runtime.GOOS).
		Str("version", runtime.Version()).
		Msg("Runtime information")

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Set default values
	defaultConfig()

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found; ignore error if desired
		var notfound viper.ConfigFileNotFoundError
		if ok := errors.Is(err, notfound); !ok {
			log.Error().Msgf("Read the config file: %s", err)
		}
	}

	viper.AutomaticEnv()

	// Replace env key
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	log.Info().
		Str("app-name", viper.GetString("APP_NAME")).
		Str("app-version", viper.GetString("APP_VERSION")).
		Str("app-port", viper.GetString("APP_PORT")).
		Msg("Config loaded")
}
