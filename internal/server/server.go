package server

import (
	"context"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/common/servercommon"
)

type server servercommon.Server

// NewServer create new server
func NewServer(commonSrv *servercommon.Server) (*server, error) {
	srv := (*server)(commonSrv)

	grpcPort := fmt.Sprintf(":%d", viper.GetInt("GRPC_PORT"))
	srv.Log.Info(fmt.Sprintf("Serving gRPC on http://localhost%s", grpcPort))

	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return nil, fmt.Errorf("net.Listen(): %w", err)
	}

	go func() {
		if err = srv.createServer(listener); err != nil {
			srv.Log.Fatal("srv.createServer()", zap.Error(err))
		}
	}()

	go func() {
		if err = srv.createGateway(grpcPort); err != nil {
			srv.Log.Fatal("srv.createGateway()", zap.Error(err))
		}
	}()

	return srv, nil
}

func (srv *server) Close() error {
	srv.GRPCServer.GracefulStop()

	return nil
}

// CreateServer create new server
func (srv *server) createServer(listener net.Listener) error {
	// Log gRPC library internals with log
	grpc_zap.ReplaceGrpcLoggerV2(srv.Log)

	streamChain := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(srv.Log),
		grpc_recovery.StreamServerInterceptor(),
		// Customer Interceptor
		srv.AuthInterceptorStream(),
	}

	unaryChain := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(srv.Log),
		grpc_recovery.UnaryServerInterceptor(),
		// Customer Interceptor
		srv.AuthInterceptorUnary(),
	}

	// Log payload if enabled
	if viper.GetBool("LOG_PAYLOAD") {
		alwaysLoggingDeciderServer := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return true
		}

		streamChain = append(streamChain, grpc_zap.PayloadStreamServerInterceptor(srv.Log, alwaysLoggingDeciderServer))
		unaryChain = append(unaryChain, grpc_zap.PayloadUnaryServerInterceptor(srv.Log, alwaysLoggingDeciderServer))
	}

	// register grpc service server
	srv.GRPCServer = grpc.NewServer(
		grpc_middleware.WithStreamServerChain(streamChain...),
		grpc_middleware.WithUnaryServerChain(unaryChain...),
	)

	// Create new validator
	// valid := validator.NewValidator(srv.Log, srv.Persist)
	// Create new validator
	// service.NewService(opts, srv.GRPCServer, valid, srv.Persist)

	if err := srv.GRPCServer.Serve(listener); err != nil {
		srv.Log.Error("srv.grpcServer.Serve()", zap.Error(err))
		return err
	}

	return nil
}
