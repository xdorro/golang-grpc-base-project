package server

import (
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/auth"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/permission"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/role"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/user"
)

func (srv *Server) registerServiceHandlers(grpcPort string, mux *runtime.ServeMux) error {
	conn, err := grpc.DialContext(
		srv.ctx,
		grpcPort,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		logger.Fatal("Failed to dial Server:", zap.Error(err))
	}

	// Register AuthService Handler
	if err = authproto.RegisterAuthServiceHandler(srv.ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterAuthServiceHandler(): %w", err)
	}

	// Register UserService Handler
	if err = userproto.RegisterUserServiceHandler(srv.ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterUserServiceHandler(): %w", err)
	}

	// Register RoleService Handler
	if err = roleproto.RegisterRoleServiceHandler(srv.ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterRoleServiceHandler(): %w", err)
	}

	// Register PermissionService Handler
	if err = permissionproto.RegisterPermissionServiceHandler(srv.ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterPermissionServiceHandler(): %w", err)
	}

	return nil
}

// createGateway create Gateway client
func (srv *Server) createGateway(grpcPort string) error {
	httpPort := fmt.Sprintf(":%d", viper.GetInt("HTTP_PORT"))
	if httpPort != "" {
		logger.Info(fmt.Sprintf("Serving gRPC-Gateway on http://localhost%s", httpPort))

		// Create HTTP Server
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

		if viper.GetBool("METRIC_ENABLE") {
			// Register Prometheus metrics handler.
			if err := mux.HandlePath("GET", "/metrics", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
				promhttp.Handler().ServeHTTP(w, r)
			}); err != nil {
				logger.Fatal("Failed to register Prometheus metrics handler:", zap.Error(err))
			}
		}

		if err := http.ListenAndServe(httpPort, mux); err != nil {
			return fmt.Errorf("http.ListenAndServe(): %w", err)
		}
	}

	return nil
}
