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

	"github.com/kucow/golang-grpc-base/internal/common"
	"github.com/kucow/golang-grpc-base/internal/service"
	"github.com/kucow/golang-grpc-base/pkg/proto/v1alpha1/helloworld"
)

// func buildDummyAuthFunction(expectedScheme string, expectedToken string) func(ctx context.Context) (
// 	context.Context, error,
// ) {
// 	return func(ctx context.Context) (context.Context, error) {
// 		token, err := grpc_auth.AuthFromMD(ctx, expectedScheme)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if token != expectedToken {
// 			return nil, status.Errorf(codes.PermissionDenied, "buildDummyAuthFunction bad token")
// 		}
// 		return context.WithValue(ctx, "some_context_marker", "marker_exists"), nil
// 	}
// }

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

	svc := service.NewService(opts)

	go func() {
		if err = srv.createServer(listener, svc); err != nil {
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

func (srv *Server) registerServiceServers(svc *service.Service) {
	helloworld.RegisterGreeterServer(srv.grpcServer, svc.HelloworldService)
}

// CreateServer create new server
func (srv *Server) createServer(listener net.Listener, svc *service.Service) error {
	streamChain := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(srv.log),
		grpc_recovery.StreamServerInterceptor(),
	}

	unaryChain := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_opentracing.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(srv.log),
		grpc_recovery.UnaryServerInterceptor(),
	}

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

	srv.registerServiceServers(svc)

	if err := srv.grpcServer.Serve(listener); err != nil {
		srv.log.Error("srv.grpcServer.Serve()", zap.Error(err))
		return err
	}

	return nil
}
