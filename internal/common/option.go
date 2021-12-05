package common

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base/pkg/ent"
	"github.com/kucow/golang-grpc-base/pkg/logger"
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
