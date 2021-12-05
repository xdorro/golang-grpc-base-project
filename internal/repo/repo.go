package repo

import (
	"context"

	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base-project/pkg/ent"
)

var _ Persist = (*Repo)(nil)

// Repo struct
type Repo struct {
	Ctx    context.Context
	Log    *zap.Logger
	Client *ent.Client
}

// NewRepo create new Persist
func NewRepo(ctx context.Context, log *zap.Logger, client *ent.Client) *Repo {
	return &Repo{
		Ctx:    ctx,
		Log:    log,
		Client: client,
	}
}
