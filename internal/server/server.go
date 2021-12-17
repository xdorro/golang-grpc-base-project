package server

import (
	"context"
	"fmt"
	"net"

	"github.com/go-redis/redis/v8"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/service"
	"github.com/xdorro/golang-grpc-base-project/pkg/client"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
)

// Server struct
type Server struct {
	ctx   context.Context
	redis redis.UniversalClient

	client     *client.Client
	grpcServer *grpc.Server
	service    *service.Service
}

// NewServer create new Server
func NewServer(ctx context.Context, client *client.Client, redis redis.UniversalClient) (*Server, error) {
	srv := &Server{
		ctx:    ctx,
		client: client,
		redis:  redis,
	}

	grpcPort := fmt.Sprintf(":%d", viper.GetInt("GRPC_PORT"))
	logger.Info(fmt.Sprintf("Serving gRPC on http://localhost%s", grpcPort))

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return nil, fmt.Errorf("net.Listen(): %w", err)
	}

	go func() {
		if err = srv.createServer(listener); err != nil {
			logger.Fatal("srv.createServer()", zap.Error(err))
		}
	}()

	go func() {
		if err = srv.createGateway(grpcPort); err != nil {
			logger.Fatal("srv.createGateway()", zap.Error(err))
		}
	}()

	return srv, nil
}

func (srv *Server) Close() error {
	srv.grpcServer.GracefulStop()

	return srv.service.Close()
}

// CreateServer create new Server
func (srv *Server) createServer(listener net.Listener) error {
	// log gRPC library internals with log
	grpc_zap.ReplaceGrpcLoggerV2(logger.NewLogger())

	streamChain := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_zap.StreamServerInterceptor(logger.NewLogger()),
		grpc_recovery.StreamServerInterceptor(),
		// Customer Interceptor
		srv.AuthInterceptorStream(),
	}

	unaryChain := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(logger.NewLogger()),
		grpc_recovery.UnaryServerInterceptor(),
		// Customer Interceptor
		srv.AuthInterceptorUnary(),
	}

	// log payload if enabled
	if viper.GetBool("LOG_PAYLOAD") {
		alwaysLoggingDeciderServer := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return true
		}

		streamChain = append(streamChain, grpc_zap.PayloadStreamServerInterceptor(logger.NewLogger(), alwaysLoggingDeciderServer))
		unaryChain = append(unaryChain, grpc_zap.PayloadUnaryServerInterceptor(logger.NewLogger(), alwaysLoggingDeciderServer))
	}

	if viper.GetBool("METRIC_ENABLE") {
		streamChain = append(streamChain, grpc_prometheus.StreamServerInterceptor)
		unaryChain = append(unaryChain, grpc_prometheus.UnaryServerInterceptor)
	}

	// register grpc service Server
	srv.grpcServer = grpc.NewServer(
		grpc_middleware.WithStreamServerChain(streamChain...),
		grpc_middleware.WithUnaryServerChain(unaryChain...),
	)

	if viper.GetBool("METRIC_ENABLE") {
		// After all your registrations, make sure all the Prometheus metrics are initialized.
		grpc_prometheus.Register(srv.grpcServer)
	}

	// Create new validator
	srv.service = service.NewService(srv.client, srv.grpcServer, srv.redis)

	if err := srv.grpcServer.Serve(listener); err != nil {
		logger.Error("srv.grpcServer.Serve()", zap.Error(err))
		return err
	}

	return nil
}
