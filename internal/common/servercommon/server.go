package servercommon

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/common/optioncommon"
	"github.com/xdorro/golang-grpc-base-project/internal/persist"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// Server struct
type Server struct {
	Ctx     context.Context
	Log     *zap.Logger
	Persist persist.Persist
	Client  *ent.Client
	Redis   redis.UniversalClient

	GRPCServer *grpc.Server
}

func NewServer(opts *optioncommon.Option) *Server {
	srv := &Server{
		Ctx:    context.Background(),
		Log:    opts.Log,
		Client: opts.Client,
		Redis:  opts.Redis,
	}

	// Create new persist
	srv.Persist = repo.NewRepo(opts.Ctx, opts.Log, opts.Client)

	return srv
}
