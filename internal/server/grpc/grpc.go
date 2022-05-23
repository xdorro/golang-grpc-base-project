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

	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

type server struct {
	sync.Mutex
	grpc               *grpc.Server
	options            []grpc.ServerOption
	streamInterceptors []grpc.StreamServerInterceptor
	unaryInterceptors  []grpc.UnaryServerInterceptor
}

// NewGrpcServer returns a IServer.
func NewGrpcServer(service service.IService, register RegisterFn) IServer {
	srv := &server{}

	logger := zerolog.InterceptorLogger(log.Logger)
	optracing := opentracing.InterceptorTracer()
	srvmetrics := metrics.NewServerMetrics()

	srv.AddStreamInterceptors(
		// tags.StreamServerInterceptor(tags.WithFieldExtractor(tags.CodeGenRequestFieldExtractor)),
		tracing.StreamServerInterceptor(optracing),
		metrics.StreamServerInterceptor(srvmetrics),
		logging.StreamServerInterceptor(logger),
		recovery.StreamServerInterceptor(),
	)

	srv.AddUnaryInterceptors(
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

		srv.AddStreamInterceptors(logging.PayloadStreamServerInterceptor(logger, alwaysLoggingDeciderServer, time.RFC3339))
		srv.AddUnaryInterceptors(logging.PayloadUnaryServerInterceptor(logger, alwaysLoggingDeciderServer, time.RFC3339))

	}

	srv.AddOptions(
		WithUnaryServerInterceptors(srv.unaryInterceptors...),
		WithStreamServerInterceptors(srv.streamInterceptors...),
	)

	srv.Lock()
	defer srv.Unlock()

	srv.grpc = grpc.NewServer(srv.options...)

	register(srv.grpc, service)

	return srv
}

func (s *server) Server() *grpc.Server {
	return s.grpc
}

func (s *server) Close() {
	s.Lock()
	defer s.Unlock()

	s.grpc.GracefulStop()
}
