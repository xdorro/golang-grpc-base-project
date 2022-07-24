package interceptor

import (
	"context"

	"github.com/bufbuild/connect-go"

	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

// NewInterceptor returns a new interceptor.
func NewInterceptor() connect.UnaryInterceptorFunc {
	interceptor := connect.UnaryInterceptorFunc(
		func(next connect.UnaryFunc) connect.UnaryFunc {
			return func(ctx context.Context, request connect.AnyRequest) (
				connect.AnyResponse, error,
			) {
				logger := log.Info()
				response, err := next(ctx, request)
				if err != nil {
					logger = log.Error()
					logger.AnErr("error", err)
				} else {
					logger.Interface("response", response.Any())
				}

				logger.Str("procedure", request.Spec().Procedure).
					Interface("request", request.Any()).
					Interface("header", request.Header()).
					Msg("Log payload interceptor")
				return response, err
			}
		},
	)

	return interceptor
}
