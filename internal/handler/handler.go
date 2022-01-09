package handler

import (
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewHandler)

type Handler struct {
	log   *zap.Logger
	repo  *repo.Repo
	redis redis.UniversalClient
}

func NewHandler(
	log *zap.Logger,
	repo *repo.Repo,
	redis redis.UniversalClient,
) *Handler {
	return &Handler{
		log:   log,
		repo:  repo,
		redis: redis,
	}
}
