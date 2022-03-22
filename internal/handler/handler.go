package handler

import (
	"context"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderHandlerSet is server providers.
var ProviderHandlerSet = wire.NewSet(NewHandler)
var _ IHandler = (*Handler)(nil)

// Handler is server struct.
type Handler struct {
	ctx  context.Context
	log  *zap.Logger
	repo repo.IRepo
}

// NewHandler creates a new service.
func NewHandler(ctx context.Context, log *zap.Logger, repo repo.IRepo) IHandler {
	return &Handler{
		ctx:  ctx,
		log:  log,
		repo: repo,
	}
}
