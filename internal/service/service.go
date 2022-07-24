package service

import (
	"net/http"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/xdorro/base-project-proto/proto-gen-go/ping/v1/pingv1connect"

	"github.com/xdorro/golang-grpc-base-project/internal/interceptor"
)

var _ IService = &Service{}

// IService is the interface that must be implemented by a service.
type IService interface {
	pingv1connect.PingServiceHandler
}

// Service is a service struct.
type Service struct {
	mux      *http.ServeMux
	services []string

	// Services
	pingv1connect.UnimplementedPingServiceHandler
}

// NewService creates a new service.
func NewService(mux *http.ServeMux) {
	s := &Service{
		mux: mux,
		services: []string{
			pingv1connect.PingServiceName,
		},
	}

	// Add connect options
	opts := connect.WithOptions(
		connect.WithCompressMinBytes(1024),
		connect.WithInterceptors(
			interceptor.NewInterceptor(),
		),
	)

	// The generated constructors return a path and a plain net/http
	// handler.
	s.Handler(opts)

	// GRPC Health
	s.Health(opts)

	// GRPC Reflect
	s.Reflect(opts)
}

// Handler is a handler.
func (s *Service) Handler(opts connect.Option) {
	s.mux.Handle(pingv1connect.NewPingServiceHandler(s, opts))
}

// Health is a health handler.
func (s *Service) Health(opts connect.Option) {
	checker := grpchealth.NewStaticChecker(s.services...)

	s.mux.Handle(grpchealth.NewHandler(checker, opts))
}

// Reflect is a reflection handler.
func (s *Service) Reflect(opts connect.Option) {
	reflector := grpcreflect.NewStaticReflector(s.services...)

	s.mux.Handle(grpcreflect.NewHandlerV1(reflector, opts))
	// Many tools still expect the older version of the server reflection API, so
	// most servers should mount both handlers.
	s.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector, opts))
}
