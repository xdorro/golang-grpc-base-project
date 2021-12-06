package common

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/kataras/jwt"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/xdorro/golang-grpc-base-project/pkg/ent"
	"github.com/xdorro/golang-grpc-base-project/pkg/ent/role"
	authproto "github.com/xdorro/golang-grpc-base-project/pkg/proto/v1/auth"
)

const (
	// TokenType token type
	TokenType = "bearer"
	// UserTokenKey user token key
	UserTokenKey = "user:%v:refresh:%s" //nolint:gosec
)

var (
	// SecretKey token secret key
	SecretKey = []byte(viper.GetString("AUTH_SECRET_KEY"))
	// AccessExpire access token expire time
	AccessExpire = 15 * time.Minute
	// RefreshExpire refresh token expire time
	RefreshExpire = 15 * time.Hour
)

// CompareHashAndPassword compare password with hash
func CompareHashAndPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateFromPassword hash password
func GenerateFromPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyToken verify token
func VerifyToken(log *zap.Logger, token string) (*jwt.VerifiedToken, error) {
	verifiedToken, err := jwt.Verify(jwt.HS256, SecretKey, []byte(token))
	if err != nil {
		log.Error("jwt.Verify()", zap.Error(err))
		return nil, TokenInvalid.Err()
	}

	return verifiedToken, nil
}

// GenerateAccessToken generate access token
func GenerateAccessToken(ctx context.Context, log *zap.Logger, user *ent.User, result *authproto.TokenResponse) error {
	now := time.Now()
	expire := now.Add(AccessExpire).Unix()
	result.AccessExpire = expire

	roles := make([]string, 0)
	perRoles, _ := user.QueryRoles().Where(role.DeleteTimeIsNil()).All(ctx)
	for _, perRole := range perRoles {
		roles = append(roles, perRole.Slug)
	}

	token, err := jwt.Sign(jwt.HS256, SecretKey, jwt.Claims{
		IssuedAt: now.Unix(),
		Expiry:   expire,
		Subject:  cast.ToString(user.ID),
		Audience: roles,
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
	expire := now.Add(RefreshExpire).Unix()
	result.RefreshExpire = expire
	token, err := jwt.Sign(jwt.HS256, SecretKey, jwt.Claims{
		IssuedAt: now.Unix(),
		Expiry:   expire,
		Subject:  cast.ToString(user.ID),
	})

	if err != nil {
		log.Error("jwt.Sign()", zap.Error(err))
		return err
	}

	result.RefreshToken = uuid.New().String()

	tokenKey := fmt.Sprintf(UserTokenKey, user.ID, result.RefreshToken)
	err = SetRefreshToken(log, rdb, tokenKey, string(token))
	if err != nil {
		return err
	}

	return nil
}

// SetRefreshToken set refresh token
func SetRefreshToken(log *zap.Logger, redis redis.UniversalClient, tokenKey, value string) error {
	err := redis.Set(redis.Context(), tokenKey, value, RefreshExpire).Err()
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
