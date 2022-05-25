package http

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

// NewHttpServer returns a IServer.
func NewHttpServer(address string, register RegisterFn) *runtime.ServeMux {
	srv := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				Multiline:       false,
				Indent:          "",
				AllowPartial:    false,
				UseProtoNames:   true,
				UseEnumNumbers:  false,
				EmitUnpopulated: false,
				Resolver:        nil,
			},
		}),
		runtime.WithForwardResponseOption(CustomForwardResponse),
		runtime.WithErrorHandler(CustomErrorResponse),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// log payload if enabled
	if viper.GetBool("LOG_PAYLOAD") {
		logger := zerolog.InterceptorLogger(log.Logger)
		alwaysLoggingDeciderClient := func(context.Context, string) logging.PayloadDecision {
			return logging.LogPayloadRequestAndResponse
		}

		opts = append(opts,
			grpc.WithUnaryInterceptor(logging.PayloadUnaryClientInterceptor(logger, alwaysLoggingDeciderClient, time.RFC3339)),
			grpc.WithStreamInterceptor(logging.PayloadStreamClientInterceptor(logger, alwaysLoggingDeciderClient, time.RFC3339)),
		)
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Panicf("Failed to dial: %s", err)
	}

	register(srv, conn)

	return srv
}
