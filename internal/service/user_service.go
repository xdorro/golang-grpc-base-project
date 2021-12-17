package service

import (
	"context"

	"github.com/spf13/cast"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/status"

	"github.com/xdorro/golang-grpc-base-project/ent"
	"github.com/xdorro/golang-grpc-base-project/pkg/common"
	"github.com/xdorro/golang-grpc-base-project/pkg/logger"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/common"
	userproto "github.com/xdorro/golang-grpc-base-project/proto/v1/user"
)

// FindAllUsers find all users
func (svc *Service) FindAllUsers(context.Context, *userproto.FindAllUsersRequest) (
	*userproto.ListUsersResponse, error,
) {
	users := svc.client.Persist.FindAllUsers()

	return &userproto.ListUsersResponse{Data: common.UsersProto(users)}, nil
}

// FindUserByID find user by id
func (svc *Service) FindUserByID(_ context.Context, in *commonproto.UUIDRequest) (
	*userproto.User, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		logger.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	u, err := svc.client.Persist.FindUserByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.UserNotExist.Err()
	}

	return common.UserProto(u), nil
}

// CreateUser create user
func (svc *Service) CreateUser(_ context.Context, in *userproto.CreateUserRequest) (
	*status.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCreateUserRequest(in); err != nil {
		logger.Error("svc.validateCreateUserRequest()", zap.Error(err))
		return nil, err
	}

	if check := svc.client.Persist.ExistUserByEmail(in.GetEmail()); check {
		return nil, common.EmailAlreadyExist.Err()
	}

	hashPassword, err := common.GenerateFromPassword(in.GetPassword())
	if err != nil {
		logger.Error("util.HashPassword()", zap.Error(err))
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
		logger.Error("svc.validateListRoles()", zap.Error(err))
		return nil, err
	}

	if err = svc.client.Persist.CreateUser(userIn, roles); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// UpdateUser update user
func (svc *Service) UpdateUser(_ context.Context, in *userproto.UpdateUserRequest) (
	*status.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateUpdateUserRequest(in); err != nil {
		logger.Error("svc.validateUpdateUserRequest()", zap.Error(err))
		return nil, err
	}

	u, err := svc.client.Persist.FindUserByID(cast.ToUint64(in.GetId()))
	if err != nil {
		return nil, common.UserNotExist.Err()
	}

	roles, err := svc.validator.ValidateListRoles(in.GetRoles())
	if err != nil {
		logger.Error("svc.validateListRoles()", zap.Error(err))
		return nil, err
	}

	u.Name = in.GetName()
	u.Email = in.GetEmail()
	u.Status = in.GetStatus()

	if err = svc.client.Persist.UpdateUser(u, roles); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}

// DeleteUser delete user
func (svc *Service) DeleteUser(_ context.Context, in *commonproto.UUIDRequest) (
	*status.Status, error,
) {
	// Validate request
	if err := svc.validator.ValidateCommonID(in); err != nil {
		logger.Error("common.ValidateCommonID()", zap.Error(err))
		return nil, err
	}

	id := cast.ToUint64(in.GetId())

	if exist := svc.client.Persist.ExistUserByID(id); !exist {
		return nil, common.UserNotExist.Err()
	}

	if err := svc.client.Persist.SoftDeleteUser(id); err != nil {
		return nil, err
	}

	return common.Success.Proto(), nil
}
