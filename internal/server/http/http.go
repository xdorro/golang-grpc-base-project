package http

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/providers/zerolog/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/xdorro/golang-grpc-base-project/internal/log"
)

type server struct {
	sync.Mutex
	address        string
	tlsCredentials credentials.TransportCredentials
	ctx            context.Context
}

// NewHttpServer returns a Server.
func NewHttpServer(address string, tlsCredentials credentials.TransportCredentials) Server {
	return &server{
		ctx:            context.Background(),
		tlsCredentials: tlsCredentials,
		address:        address,
	}
}

func (s *server) Start(register RegisterFn) *runtime.ServeMux {
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
		runtime.WithForwardResponseOption(s.CustomForwardResponse),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(s.tlsCredentials),
	}

	// log payload if enabled
	if viper.GetBool("LOG_PAYLOAD") {
		logger := zerolog.InterceptorLogger(log.Logger)
		alwaysLoggingDeciderClient := func(context.Context, string) logging.PayloadDecision {
			return logging.LogPayloadRequestAndResponse
		}

		opts = append(opts,
			grpc.WithUnaryInterceptor(logging.PayloadUnaryClientInterceptor(logger, alwaysLoggingDeciderClient, time.RFC3339)),
			grpc.WithStreamInterceptor(logging.PayloadStreamClientInterceptor(logger, alwaysLoggingDeciderClient, time.RFC3339)),
		)
	}

	conn, err := grpc.Dial(s.address, opts...)
	if err != nil {
		log.Panicf("Failed to dial: %s", err)
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				log.Errorf("Failed to close conn to %s: %v", s.address, cerr)
			}
			return
		}
		go func() {
			<-s.ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				log.Errorf("Failed to close conn to %s: %v", s.address, cerr)
			}
		}()
	}()

	s.Lock()
	register(srv, conn)
	s.Unlock()

	return srv
}

// CustomForwardResponse forwards the response from the backend to the client.
func (s *server) CustomForwardResponse(_ context.Context, w http.ResponseWriter, _ proto.Message) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])
		w.WriteHeader(http.StatusFound)
	}

	return nil
}
