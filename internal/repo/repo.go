package repo

import (
	"context"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/persist"
)

var _ persist.Persist = (*Repo)(nil)

// Repo is a wrapper around an ent.Client that provides a convenient API for
type Repo struct {
	ctx    context.Context
	client *ent.Client
}

// NewRepo create new persist
func NewRepo(ctx context.Context, client *ent.Client) *Repo {
	return &Repo{
		ctx:    ctx,
		client: client,
	}
}
