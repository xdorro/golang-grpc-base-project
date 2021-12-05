package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/kucow/golang-grpc-base/pkg/proto/v1alpha1/helloworld"
)

func (srv *Server) registerServiceHandlers(grpcPort string, mux *runtime.ServeMux) error {
	conn, err := grpc.DialContext(
		srv.ctx,
		grpcPort,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// Register Greeter
	if err = helloworld.RegisterGreeterHandler(srv.ctx, mux, conn); err != nil {
		return fmt.Errorf("helloworld.RegisterGreeterHandler(): %w", err)
	}

	return nil
}

// createGateway create Gateway client
func (srv *Server) createGateway(grpcPort string) error {
	httpPort := fmt.Sprintf(":%d", viper.GetInt("HTTP_PORT"))
	if httpPort != "" {
		srv.log.Info(fmt.Sprintf("Serving gRPC-Gateway on http://localhost%s", httpPort))

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
		if err := srv.registerServiceHandlers(grpcPort, mux); err != nil {
			return fmt.Errorf("srv.registerServiceHandlers(): %w", err)
		}

		if err := http.ListenAndServe(httpPort, mux); err != nil {
			return fmt.Errorf("http.ListenAndServe(): %w", err)
		}
	}

	return nil
}
