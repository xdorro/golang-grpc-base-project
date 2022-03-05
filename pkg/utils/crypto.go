package utils

import (
	"encoding/hex"
	"time"

	"github.com/spf13/viper"
	"github.com/vk-rv/pvx"
	"golang.org/x/crypto/bcrypt"
)

var (
	// SecretKey token secret key
	SecretKey = []byte(viper.GetString("AUTH_SECRET_KEY"))
	// AccessExpire access token expire time
	AccessExpire = 1 * time.Hour // 1 hour
	// RefreshExpire refresh token expire time
	RefreshExpire = 1 * 24 * time.Hour // 1 day
)

// GenerateFromPassword generate hash from password
func GenerateFromPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CompareHashAndPassword compare hash and password
func CompareHashAndPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}

	return true
}

// SymmetricKey create new symmetric key
func SymmetricKey() (*pvx.SymKey, error) {
	k, err := hex.DecodeString("707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f")
	if err != nil {
		return nil, err
	}

	return pvx.NewSymmetricKey(k, pvx.Version4), nil
}

// EncryptToken encrypt token
func EncryptToken(claims *pvx.RegisteredClaims) (string, error) {
	symK, err := SymmetricKey()
	if err != nil {
		return "", err
	}

	pv4 := pvx.NewPV4Local()
	token, err := pv4.Encrypt(symK, claims, pvx.WithAssert(SecretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

// DecryptToken decrypt token
func DecryptToken(token string) (*pvx.RegisteredClaims, error) {
	symK, err := SymmetricKey()
	if err != nil {
		return nil, err
	}

	pv4 := pvx.NewPV4Local()
	cc := &pvx.RegisteredClaims{}
	err = pv4.
		Decrypt(token, symK, pvx.WithAssert(SecretKey)).
		ScanClaims(cc)
	if err != nil {
		return nil, err
	}

	return cc, nil
}
