package authservice

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/persist"
	"github.com/xdorro/golang-grpc-base-project/pkg/validator"
	authproto2 "github.com/xdorro/golang-grpc-base-project/proto/v1/auth"
)

type AuthService struct {
	authproto2.UnimplementedAuthServiceServer

	ctx       context.Context
	log       *zap.Logger
	redis     redis.UniversalClient
	persist   persist.Persist
	validator *validator.Validator
}

func NewAuthService(opts *option.Option, validator *validator.Validator, persist persist.Persist) *AuthService {
	svc := &AuthService{
		ctx:       opts.Ctx,
		log:       opts.Log,
		redis:     opts.Redis,
		persist:   persist,
		validator: validator,
	}

	return svc
}

// Login login
func (svc *AuthService) Login(_ context.Context, in *authproto2.LoginRequest) (
	*authproto2.TokenResponse, error,
) {
	// Validate request
	if err := svc.validator.ValidateLoginRequest(in); err != nil {
		svc.log.Error("svc.validateLoginRequest()", zap.Error(err))
		return nil, err
	}

	u, err := svc.persist.FindUserByEmail(in.GetEmail())
	if err != nil {
		return nil, common.EmailNotExist.Err()
	}

	if !common.CompareHashAndPassword(u.Password, in.GetPassword()) {
		err = common.PasswordIncorrect.Err()
		svc.log.Error("util.CheckPasswordHash()", zap.Error(err))
		return nil, err
	}

	return svc.generateToken(u)
}

// RevokeToken revoke token
func (svc *AuthService) RevokeToken(_ context.Context, in *authproto2.TokenRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateTokenRequest(in); err != nil {
		svc.log.Error("svc.validateTokenRequest()", zap.Error(err))
		return nil, err
	}

	tokenKey := fmt.Sprintf(common.UserTokenKey, "*", in.Token)
	key, err := common.FindRefreshToken(svc.log, svc.redis, tokenKey)
	if key == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	token, err := common.GetRefreshToken(svc.log, svc.redis, key)
	if token == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	_, err = svc.validator.ValidateToken(token)
	if err != nil {
		return nil, common.TokenInvalid.Err()
	}

	if err = common.RevokeRefreshToken(svc.log, svc.redis, key); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// RefreshToken refresh token
func (svc *AuthService) RefreshToken(_ context.Context, in *authproto2.TokenRequest) (
	*authproto2.TokenResponse, error,
) {
	// Validate request
	if err := svc.validator.ValidateTokenRequest(in); err != nil {
		svc.log.Error("svc.validateTokenRequest()", zap.Error(err))
		return nil, err
	}

	tokenKey := fmt.Sprintf(common.UserTokenKey, "*", in.Token)
	key, err := common.FindRefreshToken(svc.log, svc.redis, tokenKey)
	if key == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	token, err := common.GetRefreshToken(svc.log, svc.redis, key)
	if token == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	u, err := svc.validator.ValidateToken(token)
	if err != nil {
		return nil, common.TokenInvalid.Err()
	}

	if err = common.RevokeRefreshToken(svc.log, svc.redis, key); err != nil {
		return nil, err
	}

	return svc.generateToken(u)
}

// generateToken generate token
func (svc *AuthService) generateToken(user *ent.User) (*authproto2.TokenResponse, error) {
	result := &authproto2.TokenResponse{
		TokenType: common.TokenType,
	}

	if err := common.GenerateAccessToken(svc.ctx, svc.log, user, result); err != nil {
		return nil, err
	}

	if err := common.GenerateRefreshToken(svc.log, svc.redis, user, result); err != nil {
		return nil, err
	}

	return result, nil
}
