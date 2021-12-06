package userservice

import (
	"context"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/status"

	"github.com/xdorro/golang-grpc-base-project/internal/common"
	"github.com/xdorro/golang-grpc-base-project/internal/repo"
	"github.com/xdorro/golang-grpc-base-project/internal/validator"
	"github.com/xdorro/golang-grpc-base-project/pkg/ent"
	commonproto "github.com/xdorro/golang-grpc-base-project/pkg/proto/v1/common"
	userproto "github.com/xdorro/golang-grpc-base-project/pkg/proto/v1/user"
)

type UserService struct {
	userproto.UnimplementedUserServiceServer

	log       *zap.Logger
	persist   repo.Persist
	validator *validator.Validator
}

func NewUserService(opts *common.Option, validator *validator.Validator, persist repo.Persist) *UserService {
	return &UserService{
		log:       opts.Log,
		persist:   persist,
		validator: validator,
	}
}

// FindAllUsers find all users
func (svc *UserService) FindAllUsers(context.Context, *userproto.FindAllUsersRequest) (
	*userproto.ListUsersResponse, error,
) {
	users := svc.persist.FindAllUsers()

	return &userproto.ListUsersResponse{Data: common.UsersProto(users)}, nil
}

// FindUserByID find user by id
func (svc *UserService) FindUserByID(_ context.Context, in *commonproto.UUIDRequest) (
	*userproto.User, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.persist.FindUserByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.UserNotExist.Err()
	}

	return common.UserProto(u), nil
}

// CreateUser handler CreateUser function
func (svc *UserService) CreateUser(_ context.Context, in *userproto.CreateUserRequest) (
	*status.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCreateUserRequest(in); err != nil {
		svc.log.Error("svc.validateCreateUserRequest()", zap.Error(err))
		return nil, err
	}

	if check := svc.persist.ExistUserByEmail(in.GetEmail()); check {
		return nil, common.EmailAlreadyExist.Err()
	}

	hashPassword, err := common.GenerateFromPassword(in.GetPassword())
	if err != nil {
		svc.log.Error("util.HashPassword()", zap.Error(err))
		return nil, err
	}

	userIn := &ent.User{
		Name:     in.Name,
		Email:    in.Email,
		Status:   in.Status,
		Password: hashPassword,
	}

	roles, err := svc.validator.ValidateListRoles(in.GetRoles())
	if err != nil {
		svc.log.Error("svc.validateListRoles()", zap.Error(err))
		return nil, err
	}

	if err = svc.persist.CreateUser(userIn, roles); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// UpdateUser handler UpdateUser function
func (svc *UserService) UpdateUser(_ context.Context, in *userproto.UpdateUserRequest) (
	*status.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateUpdateUserRequest(in); err != nil {
		svc.log.Error("svc.validateUpdateUserRequest()", zap.Error(err))
		return nil, err
	}

	u, err := svc.persist.FindUserByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.UserNotExist.Err()
	}

	roles, err := svc.validator.ValidateListRoles(in.GetRoles())
	if err != nil {
		svc.log.Error("svc.validateListRoles()", zap.Error(err))
		return nil, err
	}

	u.Name = in.GetName()
	u.Email = in.GetEmail()
	u.Status = in.GetStatus()

	if err = svc.persist.UpdateUser(u, roles); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// DeleteUser handler DeleteUser function
func (svc *UserService) DeleteUser(_ context.Context, in *commonproto.UUIDRequest) (
	*status.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		svc.log.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())

	if exist := svc.persist.ExistUserByID(id); !exist {
		return nil, common.UserNotExist.Err()
	}

	if err := svc.persist.SoftDeleteUser(id); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}
