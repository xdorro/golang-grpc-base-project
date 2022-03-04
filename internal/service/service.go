package service

import (
	"github.com/google/wire"
	authpb "github.com/xdorro/base-project-proto/protos/v1/auth"
	userpb "github.com/xdorro/base-project-proto/protos/v1/user"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/internal/handler"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderServiceSet is service providers.
var ProviderServiceSet = wire.NewSet(NewService)
var _ IService = (*Service)(nil)

// IService is the interface for the service
type IService interface {
	RegisterServiceServer(grpcServer *grpc.Server)
}

// Service is service struct.
type Service struct {
	log     *zap.Logger
	repo    repo.IRepo
	handler handler.IHandler

	userpb.UnimplementedUserServiceServer
	authpb.UnimplementedAuthServiceServer
}

// NewService creates a new service.
func NewService(log *zap.Logger, repo repo.IRepo, handler handler.IHandler) IService {
	return &Service{
		log:     log,
		repo:    repo,
		handler: handler,
	}
}

// RegisterServiceServer registers service server.
func (s *Service) RegisterServiceServer(grpcServer *grpc.Server) {
	userpb.RegisterUserServiceServer(grpcServer, s)
	authpb.RegisterAuthServiceServer(grpcServer, s)
}
