package helloworldservice

import (
	"context"

	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/pkg/proto/v1alpha1/helloworld"
)

type HelloworldService struct {
	helloworld.UnimplementedGreeterServer

	log *zap.Logger
	// persist repo.Persist
}

func NewHelloworldService(opts *common.Option) *HelloworldService {
	svc := &HelloworldService{
		log: opts.Log,
	}

	return svc
}

func (svc *HelloworldService) SayHello(_ context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: in.Name + " world"}, nil
}
