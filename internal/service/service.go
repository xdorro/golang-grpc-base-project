package service

import (
	"net/http"

	"github.com/bufbuild/connect-go"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"github.com/xdorro/base-project-proto/proto-gen-go/user/v1/userv1connect"
)

type IService interface {
	userv1connect.UserServiceHandler
}

// Service is a service struct.
type Service struct {
	mux      *http.ServeMux
	services []string

	// Services
	userv1connect.UnimplementedUserServiceHandler
}

// NewService creates a new service.
func NewService(mux *http.ServeMux) IService {
	s := &Service{
		mux: mux,
		services: []string{
			userv1connect.UserServiceName,
		},
	}

	// Add connect options
	opts := connect.WithOptions(
		connect.WithCompressMinBytes(1024),
	)

	// The generated constructors return a path and a plain net/http
	// handler.
	s.Handler(opts)

	// GRPC Health
	s.Health(opts)

	// GRPC Reflect
	s.Reflect(opts)

	return s
}

// Handler is a handler.
func (s *Service) Handler(opts connect.Option) {
	s.mux.Handle(userv1connect.NewUserServiceHandler(s, opts))
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
