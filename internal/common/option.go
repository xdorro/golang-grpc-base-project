package common

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base-project/internal/validator"
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	"github.com/kucow/golang-grpc-base-project/pkg/logger"
)

type Option struct {
	Ctx context.Context

	Log       *zap.Logger
	Client    *ent.Client
	Redis     redis.UniversalClient
	Validator *validator.Validator
}

func NewOption(ctx context.Context) *Option {
	return &Option{
		Ctx: ctx,
		Log: logger.NewLogger(),
	}
}
