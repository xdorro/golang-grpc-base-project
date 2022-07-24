package service

import (
	"net/http"

	"github.com/xdorro/base-project-proto/proto-gen-go/user/v1/userv1connect"
)

// Service is a service struct.
type Service struct {
	userv1connect.UnimplementedUserServiceHandler
}

// NewService creates a new service.
func NewService(mux *http.ServeMux) *Service {
	svc := &Service{}

	// The generated constructors return a path and a plain net/http
	// handler.
	mux.Handle(userv1connect.NewUserServiceHandler(svc))

	return svc
}
