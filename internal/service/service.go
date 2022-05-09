package service

import (
	"github.com/google/wire"
	"google.golang.org/grpc"

	authpb "github.com/xdorro/proto-base-project/protos/v1/auth"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"

	"github.com/xdorro/golang-grpc-base-project/internal/handler"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderServiceSet is service providers.
var ProviderServiceSet = wire.NewSet(NewService)
var _ IService = (*Service)(nil)

// Service is service struct.
type Service struct {
	repo    repo.IRepo
	handler handler.IHandler

	userpb.UnimplementedUserServiceServer
	authpb.UnimplementedAuthServiceServer
}

// NewService creates a new service.
func NewService(repo repo.IRepo, handler handler.IHandler) IService {
	return &Service{
		repo:    repo,
		handler: handler,
	}
}

// RegisterServiceServer registers service server.
func (s *Service) RegisterServiceServer(grpcServer *grpc.Server) {
	userpb.RegisterUserServiceServer(grpcServer, s)
	authpb.RegisterAuthServiceServer(grpcServer, s)
}
