package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/auth"
)

func (val *Validator) ValidateLoginRequest(in *authproto.LoginRequest) error {
	err := validation.ValidateStruct(in,
		// Validate email
		validation.Field(&in.Email,
			validation.Required,
			is.Email,
			validation.Length(5, 0),
		),
		// Validate password
		validation.Field(&in.Password,
			validation.Required,
			validation.Length(5, 0),
		),
	)

	return ValidateError(err)
}

func (val *Validator) ValidateTokenRequest(in *authproto.TokenRequest) error {
	err := validation.ValidateStruct(in,
		// Validate token
		validation.Field(&in.Token,
			validation.Required,
			validation.Length(5, 0),
		),
	)

	return ValidateError(err)
}

// ValidateToken validate token
func (val *Validator) ValidateToken(token string) (*ent.User, error) {
	verifiedToken, err := common.VerifyToken(val.log, token)
	if err != nil {
		return nil, err
	}

	userID := cast.ToUint64(verifiedToken.StandardClaims.Subject)
	u, err := val.client.Persist.FindUserByID(userID)
	if err != nil {
		err = common.UserNotExist.Err()
		val.log.Error("persist.FindUserByID()", zap.Error(err))
		return nil, err
	}

	// tokenKey := fmt.Sprintf(common.UserTokenKey, userID, token)
	// if _, err = FindRefreshToken(log, rdb, tokenKey); err != nil {
	// 	err = status.Error(codes.InvalidArgument, "Token is invalid")
	// 	return nil, err
	// }

	return u, nil
}
