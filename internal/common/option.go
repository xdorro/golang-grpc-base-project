package common

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

type Option struct {
	Ctx context.Context

	Log    *zap.Logger
	Client *ent.Client
	Redis  redis.UniversalClient
}

func NewOption(ctx context.Context) *Option {
	return &Option{
		Ctx: ctx,
		Log: logger.NewLogger(),
	}
}
