package common

import (
	"context"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// TokenType token type
	TokenType = "bearer"
	// UserSessionKey user session key
	UserSessionKey = "user:%v:session:%s" //nolint:gosec
)

var (
	// SecretKey token secret key
	SecretKey = []byte(viper.GetString("AUTH_SECRET_KEY"))
	// AccessExpire access token expire time
	AccessExpire = 1 * time.Hour // 1 hour
	// RefreshExpire refresh token expire time
	RefreshExpire = 1 * 24 * time.Hour // 1 day
)

func GetUserIDFromContext(ctx context.Context) (uint64, error) {
	md, _ := metadata.FromOutgoingContext(ctx)

	userID := md.Get(CtxUserID)
	if len(userID) == 0 {
		return 0, status.Errorf(codes.InvalidArgument, "can't retrieve user id from context")
	}

	return cast.ToUint64(userID[0]), nil
}
