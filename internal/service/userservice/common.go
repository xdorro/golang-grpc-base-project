package userservice

import (
	"github.com/kucow/golang-grpc-base-project/pkg/ent"
	userproto "github.com/kucow/golang-grpc-base-project/pkg/proto/v1/user"
)

func UserProto(user *ent.User) *userproto.User {
	return &userproto.User{
		Id:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	}
}

func UsersProto(users []*ent.User) []*userproto.User {
	result := make([]*userproto.User, len(users))

	for index, user := range users {
		result[index] = UserProto(user)
	}

	return result
}
