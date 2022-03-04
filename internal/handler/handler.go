package handler

import (
	"context"

	"github.com/google/wire"
	authpb "github.com/xdorro/base-project-proto/protos/v1/auth"
	commonpb "github.com/xdorro/base-project-proto/protos/v1/common"
	userpb "github.com/xdorro/base-project-proto/protos/v1/user"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderHandlerSet is server providers.
var ProviderHandlerSet = wire.NewSet(NewHandler)
var _ IHandler = (*Handler)(nil)

// IHandler is the interface for the server
type IHandler interface {
	ValidateError(err error) error
	ValidateCommonID(req *commonpb.UUIDRequest) error
	ValidateLoginRequest(req *authpb.LoginRequest) error
	ValidateTokenRequest(req *authpb.TokenRequest) error
	ValidateCreateUserRequest(req *userpb.CreateUserRequest) error
	ValidateUpdateUserRequest(req *userpb.UpdateUserRequest) error
}

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
