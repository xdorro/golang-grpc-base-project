package user_service

import (
	"context"

	"github.com/google/wire"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	status_proto "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	common_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/common"
	user_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/user"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/handler"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
	"github.com/xdorro/golang-grpc-base-project/internal/validator"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService)

type UserService struct {
	log       *zap.Logger
	repo      *repo.Repo
	validator *validator.Validator
	handler   *handler.Handler

	// implement UserService
	user_proto.UnimplementedUserServiceServer
}

// NewService returns a new service instance
func NewService(
	log *zap.Logger, repo *repo.Repo, validator *validator.Validator, handler *handler.Handler, grpc *grpc.Server,
) *UserService {
	svc := &UserService{
		log:       log,
		repo:      repo,
		validator: validator,
		handler:   handler,
	}

	// Register UserService Server
	user_proto.RegisterUserServiceServer(grpc, svc)

	return svc
}

// FindAllUsers find all users
func (svc *UserService) FindAllUsers(context.Context, *user_proto.FindAllUsersRequest) (
	*user_proto.ListUsersResponse, error,
) {
	users := svc.repo.FindAllUsers()

	return &user_proto.ListUsersResponse{Data: common.UsersProto(users)}, nil
}

// FindUserByID find user by id
func (svc *UserService) FindUserByID(_ context.Context, in *common_proto.UUIDRequest) (
	*user_proto.User, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.repo.FindUserByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.UserNotExist.Err()
	}

	return common.UserProto(u), nil
}

// CreateUser create user
func (svc *UserService) CreateUser(_ context.Context, in *user_proto.CreateUserRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCreateUserRequest(in); err != nil {
		svc.log.Error("svc.validateCreateUserRequest()", zap.Error(err))
		return nil, err
	}

	if check := svc.repo.ExistUserByEmail(in.GetEmail()); check {
		return nil, common.EmailAlreadyExist.Err()
	}

	role, err := svc.validator.ValidateRole(in.GetRole())
	if err != nil {
		svc.log.Error("svc.validateListRoles()", zap.Error(err))
		return nil, err
	}

	hashPassword, err := svc.handler.GenerateFromPassword(in.GetPassword())
	if err != nil {
		svc.log.Error("util.HashPassword()", zap.Error(err))
		return nil, err
	}

	userIn := &ent.User{
		Name:     in.Name,
		Email:    in.Email,
		Status:   in.Status,
		Password: hashPassword,
		RoleID:   role.ID,
	}
	if err = svc.repo.CreateUser(userIn); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// UpdateUser update user
func (svc *UserService) UpdateUser(_ context.Context, in *user_proto.UpdateUserRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateUpdateUserRequest(in); err != nil {
		svc.log.Error("svc.validateUpdateUserRequest()", zap.Error(err))
		return nil, err
	}

	u, err := svc.repo.FindUserByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.UserNotExist.Err()
	}

	role, err := svc.validator.ValidateRole(in.GetRole())
	if err != nil {
		svc.log.Error("svc.validateListRoles()", zap.Error(err))
		return nil, err
	}

	u.Name = in.GetName()
	u.Email = in.GetEmail()
	u.Status = in.GetStatus()
	u.RoleID = role.ID

	if err = svc.repo.UpdateUser(u); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// DeleteUser delete user
func (svc *UserService) DeleteUser(_ context.Context, in *common_proto.UUIDRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())

	if exist := svc.repo.ExistUserByID(id); !exist {
		return nil, common.UserNotExist.Err()
	}

	if err := svc.repo.SoftDeleteUser(id); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}
