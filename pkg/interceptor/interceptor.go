package interceptor

import (
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base-project/internal/common"
)

type Interceptor struct {
	log *zap.Logger
}

func NewInterceptor(opts *common.Option) *Interceptor {
	return &Interceptor{
		log: opts.Log,
	}
}
