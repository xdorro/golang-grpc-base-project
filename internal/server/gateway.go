package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/xdorro/golang-grpc-base-project/proto/v1/auth"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/permission"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/role"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/user"
)

func (srv *server) registerServiceHandlers(grpcPort string, mux *runtime.ServeMux) error {
	conn, err := grpc.DialContext(
		srv.Ctx,
		grpcPort,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// Register AuthService Handler
	if err = authproto.RegisterAuthServiceHandler(srv.Ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterAuthServiceHandler(): %w", err)
	}

	// Register UserService Handler
	if err = userproto.RegisterUserServiceHandler(srv.Ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterUserServiceHandler(): %w", err)
	}

	// Register RoleService Handler
	if err = roleproto.RegisterRoleServiceHandler(srv.Ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterRoleServiceHandler(): %w", err)
	}

	// Register PermissionService Handler
	if err = permissionproto.RegisterPermissionServiceHandler(srv.Ctx, mux, conn); err != nil {
		return fmt.Errorf("proto.RegisterPermissionServiceHandler(): %w", err)
	}

	return nil
}

// createGateway create Gateway client
func (srv *server) createGateway(grpcPort string) error {
	httpPort := fmt.Sprintf(":%d", viper.GetInt("HTTP_PORT"))
	if httpPort != "" {
		srv.Log.Info(fmt.Sprintf("Serving gRPC-Gateway on http://localhost%s", httpPort))

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
