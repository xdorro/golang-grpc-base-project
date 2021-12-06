package server

import (
	"context"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/kucow/golang-grpc-base-project/internal/common"
	"github.com/kucow/golang-grpc-base-project/internal/interceptor"
	"github.com/kucow/golang-grpc-base-project/internal/service"
)

// Server struct
type Server struct {
	ctx context.Context
	log *zap.Logger

	grpcServer *grpc.Server
}

// NewServer create new server
func NewServer(opts *common.Option) (*Server, error) {
	srv := &Server{
		ctx: opts.Ctx,
		log: opts.Log,
	}

	grpcPort := fmt.Sprintf(":%d", viper.GetInt("GRPC_PORT"))
	srv.log.Info(fmt.Sprintf("Serving gRPC on http://localhost%s", grpcPort))

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return nil, fmt.Errorf("net.Listen(): %w", err)
	}

	go func() {
		if err = srv.createServer(opts, listener); err != nil {
			opts.Log.Fatal("srv.createServer()", zap.Error(err))
		}
	}()

	go func() {
		if err = srv.createGateway(grpcPort); err != nil {
			opts.Log.Fatal("srv.createGateway()", zap.Error(err))
		}
	}()

	return srv, nil
}

func (srv *Server) Close() error {
	srv.grpcServer.GracefulStop()

	return nil
}

// CreateServer create new server
func (srv *Server) createServer(opts *common.Option, listener net.Listener) error {
	// Log gRPC library internals with log
	grpc_zap.ReplaceGrpcLoggerV2(srv.log)

	// Create new Interceptor
	inter := interceptor.NewInterceptor(opts.Log)

	streamChain := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(srv.log),
		grpc_auth.StreamServerInterceptor(inter.AuthInterceptor()),
		grpc_recovery.StreamServerInterceptor(),
	}

	unaryChain := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(srv.log),
		grpc_auth.UnaryServerInterceptor(inter.AuthInterceptor()),
		grpc_recovery.UnaryServerInterceptor(),
	}

	// Log payload if enabled
	if viper.GetBool("LOG_PAYLOAD") {
		alwaysLoggingDeciderServer := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return true
		}

		streamChain = append(streamChain, grpc_zap.PayloadStreamServerInterceptor(srv.log, alwaysLoggingDeciderServer))
		unaryChain = append(unaryChain, grpc_zap.PayloadUnaryServerInterceptor(srv.log, alwaysLoggingDeciderServer))
	}

	// register grpc service server
	srv.grpcServer = grpc.NewServer(
		grpc_middleware.WithStreamServerChain(streamChain...),
		grpc_middleware.WithUnaryServerChain(unaryChain...),
	)

	service.NewService(opts, srv.grpcServer)

	if err := srv.grpcServer.Serve(listener); err != nil {
		srv.log.Error("srv.grpcServer.Serve()", zap.Error(err))
		return err
	}

	return nil
}
