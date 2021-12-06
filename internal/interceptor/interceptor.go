package interceptor

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kucow/golang-grpc-base-project/internal/common"
)

type Interceptor struct {
	log *zap.Logger
}

func NewInterceptor(opts *common.Option) *Interceptor {
	return &Interceptor{
		log: opts.Log,
	}
}

func (inter *Interceptor) AuthInterceptor() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		if token != "123456789" {
			return nil, status.Errorf(codes.PermissionDenied, "buildDummyAuthFunction bad token")
		}

		return context.WithValue(ctx, "some_context_marker", "marker_exists"), nil
	}
}
