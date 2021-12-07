package repo

import (
	"context"

	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/persist"
)

var _ persist.Persist = (*Repo)(nil)

// Repo struct
type Repo struct {
	ctx    context.Context
	log    *zap.Logger
	client *ent.Client
}

// NewRepo create new persist
func NewRepo(ctx context.Context, log *zap.Logger, client *ent.Client) *Repo {
	return &Repo{
		ctx:    ctx,
		log:    log,
		client: client,
	}
}
