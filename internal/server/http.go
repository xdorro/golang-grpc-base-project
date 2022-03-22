package server

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/encoding/protojson"

	authpb "github.com/xdorro/base-project-proto/protos/v1/auth"
	userpb "github.com/xdorro/base-project-proto/protos/v1/user"
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
				EmitUnpopulated: false,
				Resolver:        nil,
			},
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCredentials),
	}

	conn, err := grpc.Dial(appPort, opts...)
	if err != nil {
		s.log.Panic("Failed to dial", zap.Error(err))
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", appPort, cerr)
			}
			return
		}
		go func() {
			<-s.ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", appPort, cerr)
			}
		}()
	}()

	// Register UserService Handler
	if err = userpb.RegisterUserServiceHandler(s.ctx, s.httpServer, conn); err != nil {
		s.log.Panic("proto.RegisterUserServiceHandler(): %w", zap.Error(err))
	}

	// Register AuthService Handler
	if err = authpb.RegisterAuthServiceHandler(s.ctx, s.httpServer, conn); err != nil {
		s.log.Panic("proto.RegisterAuthServiceHandler(): %w", zap.Error(err))
	}
}
