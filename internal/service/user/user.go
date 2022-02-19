package user_service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	status_proto "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	"github.com/xdorro/golang-grpc-base-project/api/proto/common"
	user_proto "github.com/xdorro/golang-grpc-base-project/api/proto/user"
	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/handler"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService)

var _ user_proto.UserServiceServer = (*UserService)(nil)

type UserService struct {
	log     *zap.Logger
	repo    *repo.Repo
	handler *handler.Handler
	redis   redis.UniversalClient

	// implement UserService
	user_proto.UnimplementedUserServiceServer
}

// NewService returns a new service instance
func NewService(
	log *zap.Logger, repo *repo.Repo, redis redis.UniversalClient, handler *handler.Handler, grpc *grpc.Server,
) user_proto.UserServiceServer {
	svc := &UserService{
		log:     log,
		repo:    repo,
		redis:   redis,
		handler: handler,
	}

	// Register UserService Server
	user_proto.RegisterUserServiceServer(grpc, svc)

	return svc
}

// FindAllUsers find all users
func (svc *UserService) FindAllUsers(ctx context.Context, _ *user_proto.FindAllUsersRequest) (
	*user_proto.ListUsersResponse, error,
) {
	result := &user_proto.ListUsersResponse{}

	// Get value from redis
	if value := svc.redis.Get(ctx, common.FindAllUsers).Val(); value != "" {
		if err := protojson.Unmarshal([]byte(value), result); err != nil {
			return nil, err
		}

		return result, nil
	}

	users := svc.repo.FindAllUsers()
	result = &user_proto.ListUsersResponse{Data: common.UsersProto(users)}

	// Cache data to redis
	data, err := protojson.Marshal(result)
	if err != nil {
		return nil, err
	}
	svc.redis.Set(ctx, common.FindAllUsers, data, 5*time.Hour)

	return result, nil
}

// FindUserByID find user by id
func (svc *UserService) FindUserByID(ctx context.Context, in *common_proto.UUIDRequest) (
	*user_proto.User, error,
) {
	// Validate request
	if err := svc.handler.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	result := &user_proto.User{}
	id := cast.ToUint64(in.GetId())

	// Get value from redis
	if value := svc.redis.Get(ctx, fmt.Sprintf(common.FindUserByID, id)).Val(); value != "" {
		if err := protojson.Unmarshal([]byte(value), result); err != nil {
			return nil, err
		}
		return result, nil
	}

	u, err := svc.repo.FindUserByID(id)
	if err != nil {
		return nil, common.UserNotExist.Err()
	}

	result = common.UserProto(u)

	// Cache data to redis
	data, err := protojson.Marshal(result)
	if err != nil {
		return nil, err
	}
	svc.redis.Set(ctx, fmt.Sprintf(common.FindUserByID, id), data, 5*time.Hour)

	return result, nil
}

// CreateUser create user
func (svc *UserService) CreateUser(_ context.Context, in *user_proto.CreateUserRequest) (
	*status_proto.Status, error,
) {
	// Validate request
	if err := svc.handler.ValidateCreateUserRequest(in); err != nil {
		svc.log.Error("svc.validateCreateUserRequest()", zap.Error(err))
		return nil, err
	}

	if check := svc.repo.ExistUserByEmail(in.GetEmail()); check {
		return nil, common.EmailAlreadyExist.Err()
	}

	role, err := svc.handler.ValidateRole(in.GetRole())
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
	if err := svc.handler.ValidateUpdateUserRequest(in); err != nil {
		svc.log.Error("svc.validateUpdateUserRequest()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())
	u, err := svc.repo.FindUserByID(id)
	if err != nil {
		return nil, common.UserNotExist.Err()
	}

	if check := svc.repo.ExistUserByIDNotAndEmail(id, in.GetEmail()); check {
		return nil, common.EmailAlreadyExist.Err()
	}

	role, err := svc.handler.ValidateRole(in.GetRole())
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
	if err := svc.handler.ValidateCommonID(in); err != nil {
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
