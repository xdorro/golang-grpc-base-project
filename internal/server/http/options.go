package http

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	authpb "github.com/xdorro/proto-base-project/protos/v1/auth"
	commonpb "github.com/xdorro/proto-base-project/protos/v1/common"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

// RegisterFn defines the method to register a server.
type RegisterFn func(*runtime.ServeMux, *grpc.ClientConn)

// CustomForwardResponse forwards the response from the backend to the client.
func CustomForwardResponse(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}

	return nil
}

// CustomErrorResponse custom error response
func CustomErrorResponse(
	ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request,
	err error,
) {
	path := r.URL.Path
	val, ok := runtime.RPCMethod(ctx)
	if !ok {
		log.Error().
			Str("path", path).
			Msgf("runtime.RPCMethod(): %v", err)
	} else {
		log.Info().
			Str("path", path).
			Msgf("runtime.RPCMethod(): %s", val)
	}

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

func RegisterHTTP(srv *runtime.ServeMux, conn *grpc.ClientConn) {
	ctx := context.Background()

	if err := userpb.RegisterUserServiceHandler(ctx, srv, conn); err != nil {
		log.Panicf("proto.RegisterUserServiceHandler(): %v", err)
	}

	if err := authpb.RegisterAuthServiceHandler(ctx, srv, conn); err != nil {
		log.Panicf("proto.RegisterUserServiceHandler(): %v", err)
	}
}
