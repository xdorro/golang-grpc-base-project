package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/vk-rv/pvx"
	authpb "github.com/xdorro/base-project-proto/protos/v1/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/internal/models"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
)

var _ authpb.AuthServiceServer = (*Service)(nil)

var (
	// accessExpire access token expire time
	accessExpire = 1 * time.Hour // 1 hour
	// refreshExpire refresh token expire time
	refreshExpire = 1 * 24 * time.Hour // 1 day
)

// Login is a gRPC handler for the Login method.
func (s *Service) Login(ctx context.Context, req *authpb.LoginRequest) (
	*authpb.TokenResponse, error,
) {
	// Validate request
	if err := s.handler.ValidateLoginRequest(req); err != nil {
		s.log.Error("svc.validateLoginRequest()", zap.Error(err))
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	filter := bson.M{"email": req.GetEmail()}
	data, err := s.repo.FindUser(filter)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to find user: %v", err)
	}

	if !utils.CompareHashAndPassword(data.Password, req.GetPassword()) {
		return nil, status.Errorf(codes.InvalidArgument, "password mismatch")
	}

	var result *authpb.TokenResponse
	result, err = s.generateAuthToken(data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// generateAuthToken generates a new auth token for the user.
func (s *Service) generateAuthToken(user *models.User) (
	*authpb.TokenResponse, error,
) {
	sessionID := uuid.NewString()
	now := time.Now()

	// Create a new accessToken
	accessClaims := &pvx.RegisteredClaims{
		Expiration: pvx.TimePtr(now.Add(accessExpire)),
		Subject:    cast.ToString(user.ID),
		TokenID:    sessionID,
	}
	accessToken, err := utils.EncryptToken(accessClaims)
	if err != nil {
		return nil, err
	}

	// Create a new refreshToken
	refreshClaims := &pvx.RegisteredClaims{
		Expiration: pvx.TimePtr(now.Add(refreshExpire)),
		Subject:    cast.ToString(user.ID),
		TokenID:    sessionID,
	}
	refreshToken, err := utils.EncryptToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	result := &authpb.TokenResponse{
		TokenType:     utils.TokenType,
		RefreshToken:  refreshToken,
		RefreshExpire: refreshClaims.Expiration.Unix(),
		AccessToken:   accessToken,
		AccessExpire:  accessClaims.Expiration.Unix(),
	}

	return result, nil
}
