package common

import (
	"time"

	"github.com/kataras/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		err = status.Error(codes.InvalidArgument, "Token is invalid")
		log.Error("jwt.Verify()", zap.Error(err))
		return nil, err
	}

	return verifiedToken, nil
}
