package server

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	commonpb "github.com/xdorro/base-project-proto/protos/v1/common"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	authpb "github.com/xdorro/base-project-proto/protos/v1/auth"
	userpb "github.com/xdorro/base-project-proto/protos/v1/user"

	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

// newHTTPServer create Gateway server
func (s *Server) newHTTPServer(tlsCredentials credentials.TransportCredentials, appPort string) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
		runtime.WithErrorHandler(CustomErrorHandler),
		runtime.WithForwardResponseOption(CustomForwardResponse),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(tlsCredentials),
	}

	conn, err := grpc.Dial(appPort, opts...)
	if err != nil {
		log.Panic("Failed to dial", zap.Error(err))
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
		log.Panic("proto.RegisterUserServiceHandler(): %w", zap.Error(err))
	}

	// Register AuthService Handler
	if err = authpb.RegisterAuthServiceHandler(s.ctx, s.httpServer, conn); err != nil {
		log.Panic("proto.RegisterAuthServiceHandler(): %w", zap.Error(err))
	}
}

// CustomErrorHandler handles the error from the backend to the client.
func CustomErrorHandler(
	ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request,
	err error,
) {
	val, ok := runtime.RPCMethod(ctx)
	if !ok {
		log.Error("runtime.RPCMethod(): %w", zap.Error(err))
	}
	log.Info("CustomHTTPResponse", zap.String("method", val))

	// return Internal when Marshal failed
	const fallback = `{"error": true, "message": "failed to marshal error message"}`

	var customStatus *runtime.HTTPStatusError
	if errors.As(err, &customStatus) {
		err = customStatus.Err
	}

	s := status.Convert(err)
	pb := s.Proto()

	w.Header().Del("Trailer")
	w.Header().Del("Transfer-Encoding")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encryptedResponse := &commonpb.CommonResponse{
		Error:   true,
		Message: pb.GetMessage(),
	}
	responseBody, merr := marshaler.Marshal(encryptedResponse)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", s, merr)
		if _, err = io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	if _, err = w.Write(responseBody); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}

// CustomForwardResponse forwards the response from the backend to the client.
func CustomForwardResponse(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}

	return nil
}
