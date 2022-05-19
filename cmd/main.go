package main

import (
	"github.com/rs/zerolog/log"

	"github.com/xdorro/golang-grpc-base-project/config"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	config.InitLogger()

	log.Info().
		Msg("hello world")
}
