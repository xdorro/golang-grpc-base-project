package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	status_proto "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	auth_proto "github.com/xdorro/golang-grpc-base-project/api/proto/auth"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/handler"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService)

var _ auth_proto.AuthServiceServer = (*AuthService)(nil)

type AuthService struct {
	log     *zap.Logger
	repo    *repo.Repo
	handler *handler.Handler
	redis   redis.UniversalClient

	// implement AuthService
	auth_proto.UnimplementedAuthServiceServer
}

// NewService returns a new service instance
func NewService(
	log *zap.Logger, repo *repo.Repo, redis redis.UniversalClient, grpc *grpc.Server,
	handler *handler.Handler,
) auth_proto.AuthServiceServer {
	svc := &AuthService{
		log:     log,
		repo:    repo,
		redis:   redis,
		handler: handler,
	}

	// Register AuthService Server
	auth_proto.RegisterAuthServiceServer(grpc, svc)

	return svc
}

// Login is a gRPC handler for the Login method.
func (svc *AuthService) Login(_ context.Context, in *auth_proto.LoginRequest) (
	*auth_proto.TokenResponse, error,
) {
	// Validate request
	if err := svc.handler.ValidateLoginRequest(in); err != nil {
		svc.log.Error("svc.validateLoginRequest()", zap.Error(err))
		return nil, err
	}

	u, err := svc.repo.FindUserByEmail(in.GetEmail())
	if err != nil {
		return nil, common.EmailNotExist.Err()
	}

	if !svc.handler.CompareHashAndPassword(u.Password, in.GetPassword()) {
		err = common.PasswordIncorrect.Err()
		return nil, err
	}

	now := time.Now()
	token, err := svc.handler.GenerateAuthToken(u, now)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// RevokeToken revoke token
func (svc *AuthService) RevokeToken(_ context.Context, req *auth_proto.TokenRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.handler.ValidateTokenRequest(req); err != nil {
		svc.log.Error("svc.validateTokenRequest()", zap.Error(err))
		return nil, err
	}

	claims, err := svc.handler.DecryptToken(req.GetToken())
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	tokenKey := fmt.Sprintf(common.UserSessionKey, claims.Subject, claims.TokenID)
	exist := svc.redis.Exists(svc.redis.Context(), tokenKey).Val()

	if exist < 1 {
		return nil, common.TokenInvalid.Err()
	}

	svc.log.Info("svc.handler.DecryptToken()",
		zap.Any("claims", claims),
	)

	if err = svc.redis.Del(svc.redis.Context(), tokenKey).Err(); err != nil {
		svc.log.Error("redis.Del()", zap.Error(err))
		return nil, err
	}

	return common.Success.Proto(), nil
}

// RefreshToken refresh token
func (svc *AuthService) RefreshToken(_ context.Context, req *auth_proto.TokenRequest) (
	*auth_proto.TokenResponse, error,
) {
	// Validate request
	if err := svc.handler.ValidateTokenRequest(req); err != nil {
		svc.log.Error("svc.validateTokenRequest()", zap.Error(err))
		return nil, err
	}

	claims, err := svc.handler.DecryptToken(req.GetToken())
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	tokenKey := fmt.Sprintf(common.UserSessionKey, claims.Subject, claims.TokenID)
	if err = svc.handler.ExistRefreshToken(tokenKey); err != nil {
		return nil, err
	}

	svc.log.Info("svc.handler.DecryptToken()",
		zap.Any("claims", claims),
	)

	if err = svc.redis.Del(svc.redis.Context(), tokenKey).Err(); err != nil {
		svc.log.Error("redis.Del()", zap.Error(err))
		return nil, err
	}

	u, err := svc.repo.FindUserByID(cast.ToUint64(claims.Subject))
	if err != nil {
		return nil, common.EmailNotExist.Err()
	}

	now := time.Now()
	token, err := svc.handler.GenerateAuthToken(u, now)
	if err != nil {
		return nil, err
	}

	return token, nil
}
