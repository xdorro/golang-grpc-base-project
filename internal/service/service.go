package service

import (
	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/internal/service/helloworldservice"
)

type Service struct {
	HelloworldService *helloworldservice.HelloworldService
}

func NewService(opts *common.Option) *Service {
	// Create new persist
	// persist := repo.NewRepo(opts)

	svc := &Service{
		HelloworldService: helloworldservice.NewHelloworldService(opts),
	}

	return svc
}
