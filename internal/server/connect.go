package server

import (
	"context"
	"fmt"
	"log"
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
	"github.com/kucow/golang-grpc-base/pkg/proto/v1alpha1/helloworld"
)

// CreateServer create new server
func (srv *Server) createServer(listener net.Listener, svc *service.Service) error {
	streamChain := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(srv.log),
		grpc_recovery.StreamServerInterceptor(),
	}

	unaryChain := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
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

	helloworld.RegisterGreeterServer(srv.grpcServer, svc.HelloworldService)

	if err := srv.grpcServer.Serve(listener); err != nil {
		srv.log.Error("srv.grpcServer.Serve()", zap.Error(err))
		return err
	}

	return nil
}

// listenClient listen Gateway client
func (srv *Server) listenClient(grpcPort string) error {
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

		conn, err := grpc.DialContext(
			context.Background(),
			grpcPort,
			grpc.WithBlock(),
			grpc.WithInsecure(),
		)
		if err != nil {
			log.Fatalln("Failed to dial server:", err)
		}

		mux := runtime.NewServeMux(opts...)

		// Register Greeter
		if err = helloworld.RegisterGreeterHandler(context.Background(), mux, conn); err != nil {
			return fmt.Errorf("helloworld.RegisterGreeterHandler(): %w", err)
		}

		srv.httpServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", httpPort),
			Handler: mux,
		}

		if err = srv.httpServer.ListenAndServe(); err != nil {
			return fmt.Errorf("http.ListenAndServe(): %w", err)
		}
	}

	return nil
}
