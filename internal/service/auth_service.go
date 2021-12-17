package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/pkg/common"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
	authproto "github.com/xdorro/golang-grpc-base-project/proto/v1/auth"
)

// Login is a gRPC handler for the Login method.
func (svc *Service) Login(ctx context.Context, in *authproto.LoginRequest) (
	*authproto.TokenResponse, error,
) {
	// Validate request
	if err := svc.validator.ValidateLoginRequest(in); err != nil {
		logger.Error("svc.validateLoginRequest()", zap.Error(err))
		return nil, err
	}

	u, err := svc.client.Persist.FindUserByEmail(in.GetEmail())
	if err != nil {
		return nil, common.EmailNotExist.Err()
	}

	if !common.CompareHashAndPassword(u.Password, in.GetPassword()) {
		err = common.PasswordIncorrect.Err()
		logger.Error("util.CheckPasswordHash()", zap.Error(err))
		return nil, err
	}

	return svc.generateToken(ctx, u)
}

// RevokeToken revoke token
func (svc *Service) RevokeToken(_ context.Context, in *authproto.TokenRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateTokenRequest(in); err != nil {
		logger.Error("svc.validateTokenRequest()", zap.Error(err))
		return nil, err
	}

	tokenKey := fmt.Sprintf(common.UserTokenKey, "*", in.Token)
	key, err := common.FindRefreshToken(svc.redis, tokenKey)
	if key == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	token, err := common.GetRefreshToken(svc.redis, key)
	if token == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	_, err = svc.validator.ValidateToken(token)
	if err != nil {
		return nil, common.TokenInvalid.Err()
	}

	if err = common.RevokeRefreshToken(svc.redis, key); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// RefreshToken refresh token
func (svc *Service) RefreshToken(ctx context.Context, in *authproto.TokenRequest) (
	*authproto.TokenResponse, error,
) {
	// Validate request
	if err := svc.validator.ValidateTokenRequest(in); err != nil {
		logger.Error("svc.validateTokenRequest()", zap.Error(err))
		return nil, err
	}

	tokenKey := fmt.Sprintf(common.UserTokenKey, "*", in.Token)
	key, err := common.FindRefreshToken(svc.redis, tokenKey)
	if key == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	token, err := common.GetRefreshToken(svc.redis, key)
	if token == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	u, err := svc.validator.ValidateToken(token)
	if err != nil {
		return nil, common.TokenInvalid.Err()
	}

	if err = common.RevokeRefreshToken(svc.redis, key); err != nil {
		return nil, err
	}

	return svc.generateToken(ctx, u)
}

// generateToken generate token
func (svc *Service) generateToken(ctx context.Context, user *ent.User) (*authproto.TokenResponse, error) {
	result := &authproto.TokenResponse{
		TokenType: common.TokenType,
	}

	if err := common.GenerateAccessToken(ctx, user, result); err != nil {
		return nil, err
	}

	if err := common.GenerateRefreshToken(svc.redis, user, result); err != nil {
		return nil, err
	}

	return result, nil
}
