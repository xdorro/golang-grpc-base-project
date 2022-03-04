package server

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/xdorro/base-project-proto/protos/v1/auth"
	userpb "github.com/xdorro/base-project-proto/protos/v1/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
)

// newHTTPServer create Gateway server
func (s *Server) newHTTPServer(tlsCredentials credentials.TransportCredentials, appPort string) {
	// Create HTTP Server
	s.httpServer = runtime.NewServeMux(
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
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCredentials),
	}

	// Register UserService Handler
	if err := userpb.RegisterUserServiceHandlerFromEndpoint(s.ctx, s.httpServer, appPort, opts); err != nil {
		s.log.Panic("proto.RegisterUserServiceHandler(): %w", zap.Error(err))
	}

	// Register AuthService Handler
	if err := authpb.RegisterAuthServiceHandlerFromEndpoint(s.ctx, s.httpServer, appPort, opts); err != nil {
		s.log.Panic("proto.RegisterAuthServiceHandler(): %w", zap.Error(err))
	}
}
