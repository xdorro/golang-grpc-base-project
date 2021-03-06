// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/xdorro/golang-grpc-base-project/internal/server"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

// Injectors from wire.go:

func initServer() server.IServer {
	iService := service.NewService()
	iServer := server.NewServer(iService)
	return iServer
}
