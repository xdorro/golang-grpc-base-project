package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	commonpb "github.com/xdorro/proto-base-project/protos/v1/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/xdorro/golang-grpc-base-project/pkg/log"
)

type server struct {
	sync.Mutex
	ctx context.Context
}

// NewHttpServer returns a IServer.
func NewHttpServer(address string, register RegisterFn) *runtime.ServeMux {
	s := &server{
		ctx: context.Background(),
	}

	srv := runtime.NewServeMux(
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
		runtime.WithForwardResponseOption(CustomForwardResponse),
		runtime.WithErrorHandler(CustomErrorResponse),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// log payload if enabled
	if viper.GetBool("LOG_PAYLOAD") {
		logger := zerolog.InterceptorLogger(log.Logger())
		alwaysLoggingDeciderClient := func(context.Context, string) logging.PayloadDecision {
			return logging.LogPayloadRequestAndResponse
		}

		opts = append(opts,
			grpc.WithUnaryInterceptor(logging.PayloadUnaryClientInterceptor(logger, alwaysLoggingDeciderClient, time.RFC3339)),
			grpc.WithStreamInterceptor(logging.PayloadStreamClientInterceptor(logger, alwaysLoggingDeciderClient, time.RFC3339)),
		)
	}

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Panicf("Failed to dial: %s", err)
	}

	s.Lock()
	defer s.Unlock()
	register(srv, conn)

	return srv
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

// CustomErrorResponse custom error response
func CustomErrorResponse(
	ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request,
	err error,
) {
	val, ok := runtime.RPCMethod(ctx)
	if !ok {
		log.Errorf("runtime.RPCMethod(): %v", err)
	}
	log.Infof("CustomErrorResponse method: %s", val)

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
