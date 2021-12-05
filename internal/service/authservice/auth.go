package authservice

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	statusproto "google.golang.org/genproto/googleapis/rpc/status"

	"github.com/kucow/golang-grpc-base-project/internal/common"
	"github.com/kucow/golang-grpc-base-project/internal/repo"
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	authproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/auth"
	"github.com/kucow/golang-grpc-base-project/pkg/validator"
)

type AuthService struct {
	authproto.UnimplementedAuthServiceServer

	log       *zap.Logger
	redis     redis.UniversalClient
	persist   repo.Persist
	validator *validator.Validator
}

func NewAuthService(opts *common.Option, persist repo.Persist) *AuthService {
	svc := &AuthService{
		log:       opts.Log,
		redis:     opts.Redis,
		persist:   persist,
		validator: opts.Validator,
	}

	return svc
}

// Login login
func (svc *AuthService) Login(_ context.Context, in *authproto.LoginRequest) (
	*authproto.TokenResponse, error,
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
func (svc *AuthService) RevokeToken(_ context.Context, in *authproto.TokenRequest) (
	*statusproto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateTokenRequest(in); err != nil {
		svc.log.Error("svc.validateTokenRequest()", zap.Error(err))
		return nil, err
	}

	tokenKey := fmt.Sprintf(common.UserTokenKey, "*", in.Token)
	key, err := FindRefreshToken(svc.log, svc.redis, tokenKey)
	if key == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	token, err := GetRefreshToken(svc.log, svc.redis, key)
	if token == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	_, err = ValidateToken(svc.log, svc.persist, token)
	if err != nil {
		return nil, common.TokenInvalid.Err()
	}

	if err = RevokeRefreshToken(svc.log, svc.redis, key); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// RefreshToken refresh token
func (svc *AuthService) RefreshToken(_ context.Context, in *authproto.TokenRequest) (
	*authproto.TokenResponse, error,
) {
	// Validate request
	if err := svc.validator.ValidateTokenRequest(in); err != nil {
		svc.log.Error("svc.validateTokenRequest()", zap.Error(err))
		return nil, err
	}

	tokenKey := fmt.Sprintf(common.UserTokenKey, "*", in.Token)
	key, err := FindRefreshToken(svc.log, svc.redis, tokenKey)
	if key == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	token, err := GetRefreshToken(svc.log, svc.redis, key)
	if token == "" || err != nil {
		return nil, common.TokenInvalid.Err()
	}

	u, err := ValidateToken(svc.log, svc.persist, token)
	if err != nil {
		return nil, common.TokenInvalid.Err()
	}

	if err = RevokeRefreshToken(svc.log, svc.redis, key); err != nil {
		return nil, err
	}

	return svc.generateToken(u)
}

// generateToken generate token
func (svc *AuthService) generateToken(user *ent.User) (*authproto.TokenResponse, error) {
	result := &authproto.TokenResponse{
		TokenType: common.TokenType,
	}

	if err := GenerateAccessToken(svc.log, user, result); err != nil {
		return nil, err
	}

	if err := GenerateRefreshToken(svc.log, svc.redis, user, result); err != nil {
		return nil, err
	}

	return result, nil
}
