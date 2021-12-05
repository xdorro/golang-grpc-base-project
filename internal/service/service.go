package service

import (
	"github.com/kucow/golang-grpc-base/internal/common"
)

type Service struct {
}

func NewService(opts *common.Option) *Service {
	// Create new persist
	// persist := repo.NewRepo(opts)

	svc := &Service{}

	return svc
}
