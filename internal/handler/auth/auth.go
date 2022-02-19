package auth_handler

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/spf13/cast"
	"github.com/vk-rv/pvx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/proto/auth"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderSet is handler providers.
var ProviderSet = wire.NewSet(NewHandler)

var _ AuthPersist = (*AuthHandler)(nil)

type AuthHandler struct {
	log   *zap.Logger
	repo  *repo.Repo
	redis redis.UniversalClient
}

func NewHandler(
	log *zap.Logger,
	repo *repo.Repo,
	redis redis.UniversalClient,
) AuthPersist {
	return &AuthHandler{
		log:   log,
		repo:  repo,
		redis: redis,
	}
}

// GenerateFromPassword hash password
func (handler *AuthHandler) GenerateFromPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CompareHashAndPassword compare password with hash
func (handler *AuthHandler) CompareHashAndPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		handler.log.Error("bcrypt.CompareHashAndPassword()", zap.Error(err))
		return false
	}

	return true
}

func (handler *AuthHandler) SymmetricKey() (*pvx.SymKey, error) {
	k, err := hex.DecodeString("707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f")
	if err != nil {
		handler.log.Error("handler.SymmetricKey()", zap.Error(err))
		return nil, err
	}

	return pvx.NewSymmetricKey(k, pvx.Version4), nil
}

func (handler *AuthHandler) EncryptToken(claims *pvx.RegisteredClaims) (string, error) {
	symK, err := handler.SymmetricKey()
	if err != nil {
		return "", err
	}

	pv4 := pvx.NewPV4Local()
	token, err := pv4.Encrypt(symK, claims, pvx.WithAssert(common.SecretKey))
	if err != nil {
		handler.log.Error("pv4.Encrypt()", zap.Error(err))
		return "", err
	}

	return token, nil
}

func (handler *AuthHandler) DecryptToken(token string) (*pvx.RegisteredClaims, error) {
	symK, err := handler.SymmetricKey()
	if err != nil {
		return nil, err
	}

	pv4 := pvx.NewPV4Local()
	cc := &pvx.RegisteredClaims{}
	err = pv4.Decrypt(token, symK, pvx.WithAssert(common.SecretKey)).
		ScanClaims(cc)
	if err != nil {
		handler.log.Error("pv4.Decrypt()", zap.Error(err))
		return nil, err
	}

	return cc, nil
}

func (handler *AuthHandler) GenerateAccessClaims(
	user *ent.User, now time.Time, sessionID string,
) *pvx.RegisteredClaims {
	claims := &pvx.RegisteredClaims{
		Expiration: pvx.TimePtr(now.Add(common.AccessExpire)),
		Subject:    cast.ToString(user.ID),
		TokenID:    sessionID,
	}

	return claims
}

func (handler *AuthHandler) GenerateRefreshClaims(
	user *ent.User, now time.Time, sessionID string,
) *pvx.RegisteredClaims {
	claims := &pvx.RegisteredClaims{
		Expiration: pvx.TimePtr(now.Add(common.RefreshExpire)),
		Subject:    cast.ToString(user.ID),
		TokenID:    sessionID,
	}

	return claims
}

func (handler *AuthHandler) GenerateAuthToken(user *ent.User, now time.Time) (*auth_proto.TokenResponse, error) {
	sessionID := uuid.NewString()

	refreshClaims := handler.GenerateRefreshClaims(user, now, sessionID)
	refreshToken, err := handler.EncryptToken(refreshClaims)
	if err != nil {
		return nil, err
	}

	accessClaims := handler.GenerateAccessClaims(user, now, sessionID)
	accessToken, err := handler.EncryptToken(accessClaims)
	if err != nil {
		return nil, err
	}

	result := &auth_proto.TokenResponse{
		TokenType:     common.TokenType,
		RefreshToken:  refreshToken,
		RefreshExpire: refreshClaims.Expiration.Unix(),
		AccessToken:   accessToken,
		AccessExpire:  accessClaims.Expiration.Unix(),
	}

	tokenKey := fmt.Sprintf(common.UserSessionKey, user.ID, refreshClaims.TokenID)
	err = handler.redis.Set(handler.redis.Context(), tokenKey, refreshToken, common.RefreshExpire).Err()
	if err != nil {
		handler.log.Error("redis.Set()", zap.Error(err))
		return nil, err
	}

	return result, nil
}

func (handler *AuthHandler) ExistRefreshToken(tokenKey string) error {
	exist := handler.redis.Exists(handler.redis.Context(), tokenKey).Val()

	if exist < 1 {
		return common.TokenInvalid.Err()
	}

	return nil
}
