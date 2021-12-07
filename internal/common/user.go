package common

import (
	"github.com/spf13/cast"

	"github.com/xdorro/golang-grpc-base-project/pkg/ent"
	"github.com/xdorro/golang-grpc-base-project/proto/v1/user"
)

// UserProto convert ent user to proto
func UserProto(user *ent.User) *userproto.User {
	return &userproto.User{
		Id:     cast.ToString(user.ID),
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	}
}

// UsersProto convert ent users to proto
func UsersProto(users []*ent.User) []*userproto.User {
	result := make([]*userproto.User, len(users))

	for index, user := range users {
		result[index] = UserProto(user)
	}

	return result
}
