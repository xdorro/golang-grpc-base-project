package main

import (
	"github.com/rs/zerolog/log"

	"github.com/xdorro/golang-grpc-base-project/config"
)

func main() {
	config.InitConfig()

	log.Info().
		Msg("hello world")
}
