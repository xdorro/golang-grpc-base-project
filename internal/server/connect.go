package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kucow/golang-grpc-base/internal/service"
)

// CreateServer create new server
func (srv *Server) createServer(listener net.Listener) error {
	// Log gRPC library internals with log
	grpc_zap.ReplaceGrpcLoggerV2(srv.log)

	chain := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(srv.log),
		grpc_recovery.UnaryServerInterceptor(),
	}

	if viper.GetBool("LOG_PAYLOAD") {
		decider := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return true
		}

		chain = append(chain, grpc_zap.PayloadUnaryServerInterceptor(srv.log, decider))
	}

	// register grpc service server
	srv.grpcServer = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(chain...),
	)

	if err := srv.grpcServer.Serve(listener); err != nil {
		srv.log.Error("srv.grpcServer.Serve()", zap.Error(err))
		return err
	}

	return nil
}

// listenServer listen gRPC server
func (srv *Server) listenServer() (net.Listener, error) {
	grpcPort := viper.GetInt("GRPC_PORT")
	srv.log.Info(fmt.Sprintf("Serving gRPC on http://localhost:%d", grpcPort))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return nil, err
	}

	return listener, nil
}

// listenClient listen Gateway client
func (srv *Server) listenClient(svc *service.Service) error {
	httpPort := viper.GetInt("HTTP_PORT")
	if httpPort != 0 {
		srv.log.Info(fmt.Sprintf("Serving gRPC-Gateway on http://localhost:%d", httpPort))

		// Create HTTP server
		opts := []runtime.ServeMuxOption{
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					Multiline:       false,
					Indent:          "",
					AllowPartial:    false,
					UseProtoNames:   true,
					UseEnumNumbers:  false,
					EmitUnpopulated: true,
					Resolver:        nil,
				},
			}),
		}

		mux := runtime.NewServeMux(opts...)

		srv.httpServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", httpPort),
			Handler: mux,
		}

		if err := srv.httpServer.ListenAndServe(); err != nil {
			return fmt.Errorf("http.ListenAndServe(): %w", err)
		}
	}

	return nil
}
