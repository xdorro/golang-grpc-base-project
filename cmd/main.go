package main

import (
	"github.com/xdorro/golang-grpc-base-project/config"
	"github.com/xdorro/golang-grpc-base-project/internal/log"
)

func main() {
	config.InitConfig()

	log.Infof("hello world")
}
