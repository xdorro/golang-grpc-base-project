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

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/internal/common/option"
	"github.com/xdorro/golang-grpc-base-project/internal/persist"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// Server struct
type Server struct {
	ctx     context.Context
	log     *zap.Logger
	Persist persist.Persist
	Client  *ent.Client
	Redis   redis.UniversalClient

	grpcServer *grpc.Server
}

// NewServer create new Server
func NewServer(opts *option.Option) (*Server, error) {
	srv := &Server{
		ctx:    opts.Ctx,
		log:    opts.Log,
		Client: opts.Client,
		Redis:  opts.Redis,
	}

	// Create new persist
	srv.Persist = repo.NewRepo(opts.Ctx, opts.Log, opts.Client)

	grpcPort := fmt.Sprintf(":%d", viper.GetInt("GRPC_PORT"))
	srv.log.Info(fmt.Sprintf("Serving gRPC on http://localhost%s", grpcPort))

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return nil, fmt.Errorf("net.Listen(): %w", err)
	}

	go func() {
		if err = srv.createServer(listener); err != nil {
			srv.log.Fatal("srv.createServer()", zap.Error(err))
		}
	}()

	go func() {
		if err = srv.createGateway(grpcPort); err != nil {
			srv.log.Fatal("srv.createGateway()", zap.Error(err))
		}
	}()

	return srv, nil
}

func (srv *Server) Close() error {
	srv.grpcServer.GracefulStop()

	return nil
}

// CreateServer create new Server
func (srv *Server) createServer(listener net.Listener) error {
	// Log gRPC library internals with log
	grpc_zap.ReplaceGrpcLoggerV2(srv.log)

	streamChain := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(srv.log),
		grpc_recovery.StreamServerInterceptor(),
		// Customer Interceptor
		srv.AuthInterceptorStream(),
	}

	unaryChain := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(srv.log),
		grpc_recovery.UnaryServerInterceptor(),
		// Customer Interceptor
		srv.AuthInterceptorUnary(),
	}

	// Log payload if enabled
	if viper.GetBool("LOG_PAYLOAD") {
		alwaysLoggingDeciderServer := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return true
		}

		streamChain = append(streamChain, grpc_zap.PayloadStreamServerInterceptor(srv.log, alwaysLoggingDeciderServer))
		unaryChain = append(unaryChain, grpc_zap.PayloadUnaryServerInterceptor(srv.log, alwaysLoggingDeciderServer))
	}

	// register grpc service Server
	srv.grpcServer = grpc.NewServer(
		grpc_middleware.WithStreamServerChain(streamChain...),
		grpc_middleware.WithUnaryServerChain(unaryChain...),
	)

	// Create new validator
	// valid := validator.NewValidator(srv.log, srv.Persist)
	// Create new validator
	// service.NewService(opts, srv.grpcServer, valid, srv.Persist)

	if err := srv.grpcServer.Serve(listener); err != nil {
		srv.log.Error("srv.grpcServer.Serve()", zap.Error(err))
		return err
	}

	return nil
}
