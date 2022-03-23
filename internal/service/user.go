package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonpb "github.com/xdorro/proto-base-project/protos/v1/common"
	userpb "github.com/xdorro/proto-base-project/protos/v1/user"

	"github.com/xdorro/golang-grpc-base-project/internal/models"
	"github.com/xdorro/golang-grpc-base-project/pkg/log"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
)

// FindAllUsers returns all users
func (s *Service) FindAllUsers(_ context.Context, req *userpb.FindAllUsersRequest) (
	*userpb.ListUsersResponse, error,
) {
	// count all users with filter
	filter := bson.M{}
	count, _ := s.repo.CountAllUsers(filter)
	limit := int64(10)
	totalPages := utils.TotalPage(count, limit)
	page := utils.CurrentPage(req.GetPage(), totalPages)

	// find all genres with filter and option
	opt := options.
		Find().
		SetSort(bson.M{"created_at": -1}).
		SetProjection(bson.M{"password": 0}).
		SetLimit(limit).
		SetSkip((page - 1) * limit)
	data, err := s.repo.FindAllUsers(filter, opt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find all users: %v", err)
	}

	result := &userpb.ListUsersResponse{
		Data:        models.ToUsersProto(data),
		TotalPage:   totalPages,
		CurrentPage: page,
	}

	return result, nil
}

// FindUserByID returns a user by id
func (s *Service) FindUserByID(_ context.Context, req *commonpb.UUIDRequest) (
	*userpb.User, error,
) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	opt := options.
		FindOne().
		SetSort(bson.M{"created_at": -1}).
		SetProjection(bson.M{"password": 0})
	filter := bson.D{{"_id", id}}

	result, err := s.repo.FindUser(filter, opt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	return result.ToProto(), nil
}

// CreateUser creates a user
func (s *Service) CreateUser(_ context.Context, req *userpb.CreateUserRequest) (
	*statuspb.Status, error,
) {
	// Validate request
	if err := s.handler.ValidateCreateUserRequest(req); err != nil {
		log.Error("svc.validateCreateUserRequest()", zap.Error(err))
		return nil, err
	}

	// count all users with filter
	count, _ := s.repo.CountAllUsers(bson.M{"email": req.GetEmail()})
	if count > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "user email already exists")
	}

	password, err := utils.GenerateFromPassword(req.GetPassword())
	if err != nil {
		log.Error("Error create user", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "failed to create user: %v", err)
	}

	data := &models.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: password,
	}
	data.BeforeCreate()

	if err = s.repo.CreateUser(data); err != nil {
		log.Error("Error creating user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return status.New(codes.OK, "create success").Proto(), nil
}

// UpdateUser updates a user
func (s *Service) UpdateUser(_ context.Context, req *userpb.UpdateUserRequest) (
	*statuspb.Status, error,
) {
	// Validate request
	if err := s.handler.ValidateUpdateUserRequest(req); err != nil {
		log.Error("svc.validateUpdateUserRequest()", zap.Error(err))
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", id}}
	data, err := s.repo.FindUser(filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	// count all users with filter
	count, _ := s.repo.CountAllUsers(bson.M{
		"_id":   bson.M{"$ne": id},
		"email": req.GetEmail(),
	})
	if count > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "user email already exists")
	}

	data.Name = utils.StringCompareOrPassValue(data.Name, req.GetName())
	data.Email = utils.StringCompareOrPassValue(data.Email, req.GetEmail())

	if err = s.repo.UpdateUser(filter, bson.M{"$set": data}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return status.New(codes.OK, "update success").Proto(), nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(_ context.Context, req *commonpb.UUIDRequest) (
	*statuspb.Status, error,
) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}
	// count all users with filter
	count, _ := s.repo.CountAllUsers(filter)
	if count <= 0 {
		return nil, status.Errorf(codes.AlreadyExists, "user does not exists")
	}

	if err = s.repo.DeleteUser(filter); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return status.New(codes.OK, "delete success").Proto(), nil
}
