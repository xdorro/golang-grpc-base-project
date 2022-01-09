package server

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

// NewHTTPServer create Gateway client
func NewHTTPServer() *runtime.ServeMux {
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

	return runtime.NewServeMux(opts...)
}
