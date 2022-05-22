package grpc

import (
	"context"
	"sync"
	"time"

	metrics "github.com/grpc-ecosystem/go-grpc-middleware/providers/openmetrics/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/providers/opentracing/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/tracing"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/log"
	"github.com/xdorro/golang-grpc-base-project/internal/service"
)

type server struct {
	sync.Mutex
	options            []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
}

// NewGrpcServer returns a Server.
func NewGrpcServer() Server {
	return &server{}
}

func (s *server) Start(register RegisterFn) *grpc.Server {
	logger := zerolog.InterceptorLogger(log.Logger)
	optracing := opentracing.InterceptorTracer()
	srvmetrics := metrics.NewServerMetrics()

	s.AddStreamInterceptors(
		// tags.StreamServerInterceptor(tags.WithFieldExtractor(tags.CodeGenRequestFieldExtractor)),
		tracing.StreamServerInterceptor(optracing),
		metrics.StreamServerInterceptor(srvmetrics),
		logging.StreamServerInterceptor(logger),
		recovery.StreamServerInterceptor(),
	)

	s.AddUnaryInterceptors(
		// tags.UnaryServerInterceptor(tags.WithFieldExtractor(tags.CodeGenRequestFieldExtractor)),
		tracing.UnaryServerInterceptor(optracing),
		metrics.UnaryServerInterceptor(srvmetrics),
		logging.UnaryServerInterceptor(logger),
		recovery.UnaryServerInterceptor(),
	)

	// log payload if enabled
	if viper.GetBool("LOG_PAYLOAD") {
		alwaysLoggingDeciderServer := func(context.Context, string, interface{}) logging.PayloadDecision {
			return logging.LogPayloadRequestAndResponse
		}

		s.AddStreamInterceptors(logging.PayloadStreamServerInterceptor(logger, alwaysLoggingDeciderServer, time.RFC3339))
		s.AddUnaryInterceptors(logging.PayloadUnaryServerInterceptor(logger, alwaysLoggingDeciderServer, time.RFC3339))

	}

	s.AddOptions(
		WithUnaryServerInterceptors(s.unaryInterceptors...),
		WithStreamServerInterceptors(s.streamInterceptors...),
	)

	srv := grpc.NewServer(s.options...)
	svc := service.NewService()
	s.Lock()
	register(srv, svc)
	s.Unlock()

	return srv
}

func (s *server) Stop() error {
	return nil
}
