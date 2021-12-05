package repo

import (
	"context"

	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base-project/internal/common"
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
func NewRepo(opts *common.Option) *Repo {
	return &Repo{
		Ctx:    opts.Ctx,
		Log:    opts.Log,
		Client: opts.Client,
	}
}
