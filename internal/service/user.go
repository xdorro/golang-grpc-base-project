package service

import (
	"context"
	"strings"
	"time"

	commonpb "github.com/xdorro/base-project-proto/protos/v1/common"
	userpb "github.com/xdorro/base-project-proto/protos/v1/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	statuspb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xdorro/golang-grpc-base-project/internal/models"
	"github.com/xdorro/golang-grpc-base-project/pkg/utils"
)

// FindAllUsers returns all users
func (s *Service) FindAllUsers(ctx context.Context, _ *userpb.FindAllUsersRequest) (
	*userpb.ListUsersResponse, error,
) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	opt := options.
		Find().
		SetSort(bson.M{"created_at": -1}).
		SetProjection(bson.M{"password": 0})

	filter := bson.D{}
	cur, err := s.repo.
		UserCollection().
		Find(ctx, filter, opt)
	if err != nil {
		s.log.Error("Error find all users", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to find all users: %v", err)
	}

	defer func() {
		_ = cur.Close(ctx)
	}()

	var data []*userpb.User
	for cur.Next(ctx) {
		user := &models.User{}
		if err = cur.Decode(user); err != nil {
			s.log.Error("Error find all users", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "failed to find all users: %v", err)
		}

		data = append(data, user.UserToProto())
	}

	result := &userpb.ListUsersResponse{
		Data: data,
	}

	return result, nil
}

// FindUserByID returns a user by id
func (s *Service) FindUserByID(ctx context.Context, req *commonpb.UUIDRequest) (
	*userpb.User, error,
) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	opt := options.
		FindOne().
		SetSort(bson.M{"created_at": -1}).
		SetProjection(bson.M{"password": 0})
	filter := bson.D{{"_id", id}}
	result := &models.User{}
	err = s.repo.
		UserCollection().
		FindOne(ctx, filter, opt).
		Decode(result)

	if err != nil {
		s.log.Error("Error find user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	s.log.Info("Find user by id", zap.Any("user", result))

	return result.UserToProto(), nil
}

// CreateUser creates a user
func (s *Service) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (
	*statuspb.Status, error,
) {
	// Validate request
	if err := s.handler.ValidateCreateUserRequest(req); err != nil {
		s.log.Error("svc.validateCreateUserRequest()", zap.Error(err))
		return nil, err
	}

	password, err := utils.GenerateFromPassword(req.GetPassword())
	if err != nil {
		s.log.Error("Error create user", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "failed to create user: %v", err)
	}

	data := &models.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: password,
	}
	data.BeforeCreate()

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	_, err = s.repo.
		UserCollection().InsertOne(ctx, data)
	if err != nil {
		s.log.Error("Error creating user", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return status.New(codes.OK, "create success").Proto(), nil
}

// UpdateUser updates a user
func (s *Service) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (
	*statuspb.Status, error,
) {
	// Validate request
	if err := s.handler.ValidateUpdateUserRequest(req); err != nil {
		s.log.Error("svc.validateUpdateUserRequest()", zap.Error(err))
		return nil, err
	}

	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	filter := bson.D{{"_id", id}}
	data := &models.User{}
	err = s.repo.
		UserCollection().
		FindOne(ctx, filter).
		Decode(data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	if strings.Compare(req.GetName(), data.Name) != 0 {
		data.Name = req.GetName()
	}

	if strings.Compare(req.GetEmail(), data.Email) != 0 {
		data.Email = req.GetEmail()
	}

	_, err = s.repo.
		UserCollection().
		UpdateOne(ctx, filter, bson.D{
			{"$set", data},
		})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return status.New(codes.OK, "update success").Proto(), nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, req *commonpb.UUIDRequest) (
	*statuspb.Status, error,
) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	filter := bson.D{{"_id", id}}
	result := s.repo.
		UserCollection().
		FindOneAndDelete(ctx, filter)
	if err = result.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return status.New(codes.OK, "delete success").Proto(), nil
}
