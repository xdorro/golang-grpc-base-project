package authservice

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/kataras/jwt"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/internal/repo"
	"github.com/kucow/golang-grpc-base/pkg/ent"
	authproto "github.com/kucow/golang-grpc-base/pkg/proto/v1/auth"
)

// ValidateToken validate token
func ValidateToken(log *zap.Logger, persist repo.Persist, token string) (*ent.User, error) {
	verifiedToken, err := common.VerifyToken(log, token)
	if err != nil {
		return nil, err
	}

	userID := cast.ToString(verifiedToken.StandardClaims.Subject)
	u, err := persist.FindUserByID(userID)
	if err != nil {
		err = common.UserNotExist.Err()
		log.Error("persist.FindUserByID()", zap.Error(err))
		return nil, err
	}

	// tokenKey := fmt.Sprintf(common.UserTokenKey, userID, token)
	// if _, err = FindRefreshToken(log, rdb, tokenKey); err != nil {
	// 	err = status.Error(codes.InvalidArgument, "Token is invalid")
	// 	return nil, err
	// }

	return u, nil
}

// GenerateAccessToken generate access token
func GenerateAccessToken(log *zap.Logger, user *ent.User, result *authproto.TokenResponse) error {
	now := time.Now()
	expire := now.Add(common.AccessExpire).Unix()
	result.AccessExpire = expire
	token, err := jwt.Sign(jwt.HS256, common.SecretKey, jwt.Claims{
		IssuedAt: now.Unix(),
		Expiry:   expire,
		Subject:  cast.ToString(user.ID),
	})

	if err != nil {
		log.Error("jwt.Sign()", zap.Error(err))
		return err
	}

	result.AccessToken = string(token)
	return nil
}

// GenerateRefreshToken generate refresh token
func GenerateRefreshToken(
	log *zap.Logger, rdb redis.UniversalClient, user *ent.User, result *authproto.TokenResponse,
) error {
	now := time.Now()
	expire := now.Add(common.RefreshExpire).Unix()
	result.RefreshExpire = expire
	token, err := jwt.Sign(jwt.HS256, common.SecretKey, jwt.Claims{
		IssuedAt: now.Unix(),
		Expiry:   expire,
		Subject:  cast.ToString(user.ID),
	})

	if err != nil {
		log.Error("jwt.Sign()", zap.Error(err))
		return err
	}

	result.RefreshToken = uuid.New().String()

	tokenKey := fmt.Sprintf(common.UserTokenKey, user.ID, result.RefreshToken)
	err = SetRefreshToken(log, rdb, tokenKey, string(token))
	if err != nil {
		return err
	}

	return nil
}

// SetRefreshToken set refresh token
func SetRefreshToken(log *zap.Logger, redis redis.UniversalClient, tokenKey, value string) error {
	err := redis.Set(redis.Context(), tokenKey, value, common.RefreshExpire).Err()
	if err != nil {
		log.Error("redis.Set()", zap.Error(err))
		return err
	}

	return nil
}

// FindRefreshToken find refresh token
func FindRefreshToken(log *zap.Logger, rdb redis.UniversalClient, tokenKey string) (string, error) {
	keys, _, err := rdb.Scan(rdb.Context(), 0, tokenKey, 0).Result()
	if err != nil {
		log.Error("redis.Keys()", zap.Error(err))
		return "", err
	}

	if len(keys) > 0 {
		return keys[0], nil
	}

	return "", nil
}

// GetRefreshToken get refresh token
func GetRefreshToken(log *zap.Logger, rdb redis.UniversalClient, tokenKey string) (string, error) {
	val, err := rdb.Get(rdb.Context(), tokenKey).Result()
	if err != nil {
		log.Error("redis.Get()", zap.Error(err))
		return "", err
	}

	return val, nil
}

// RevokeRefreshToken revoke refresh token
func RevokeRefreshToken(log *zap.Logger, rdb redis.UniversalClient, tokenKey string) error {
	err := rdb.Del(rdb.Context(), tokenKey).Err()
	if err != nil {
		log.Error("rdb.Del()", zap.Error(err))
		return err
	}

	return nil
}
