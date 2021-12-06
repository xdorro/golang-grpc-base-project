package interceptor

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/zap"

	"github.com/kucow/golang-grpc-base-project/internal/common"
)

type Interceptor struct {
	log *zap.Logger
}

func NewInterceptor(log *zap.Logger) *Interceptor {
	return &Interceptor{
		log: log,
	}
}

func (inter *Interceptor) AuthInterceptor() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, common.TokenType)
		if err != nil {
			return nil, err
		}

		_, err = common.VerifyToken(inter.log, token)
		if err != nil {
			return nil, err
		}

		// userID := cast.ToUint64(verifiedToken.StandardClaims.Subject)
		// if token != "123456789" {
		// 	return nil, status.Errorf(codes.PermissionDenied, "bad token")
		// }

		return context.WithValue(ctx, "some_context_marker", "marker_exists"), nil
	}
}
