package service

import (
	"google.golang.org/grpc"

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/internal/service/helloworldservice"
	"github.com/kucow/golang-grpc-base/pkg/proto/v1alpha1/helloworld"
)

func NewService(opts *common.Option, srv *grpc.Server) {
	// Create new persist
	// persist := repo.NewRepo(opts)

	helloworldsvc := helloworldservice.NewHelloworldService(opts)

	helloworld.RegisterGreeterServer(srv, helloworldsvc)
}
