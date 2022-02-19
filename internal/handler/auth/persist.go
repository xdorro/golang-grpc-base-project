package auth_handler

import (
	"time"

	"github.com/vk-rv/pvx"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	auth_proto "github.com/xdorro/golang-grpc-base-project/api/proto/auth"
)

type AuthPersist interface {
	GenerateFromPassword(password string) (string, error)

	CompareHashAndPassword(hash, password string) bool

	SymmetricKey() (*pvx.SymKey, error)

	EncryptToken(claims *pvx.RegisteredClaims) (string, error)

	DecryptToken(token string) (*pvx.RegisteredClaims, error)

	GenerateAccessClaims(user *ent.User, now time.Time, sessionID string) *pvx.RegisteredClaims

	GenerateRefreshClaims(user *ent.User, now time.Time, sessionID string) *pvx.RegisteredClaims

	GenerateAuthToken(user *ent.User, now time.Time) (*auth_proto.TokenResponse, error)

	ExistRefreshToken(tokenKey string) error
}
