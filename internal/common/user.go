package common

import (
	"github.com/spf13/cast"

	"github.com/xdorro/golang-grpc-base-project/api/ent"
	user_proto "github.com/xdorro/golang-grpc-base-project/api/proto/v1/user"
)

// UserProto convert ent user to proto
func UserProto(user *ent.User) *user_proto.User {
	return &user_proto.User{
		Id:     cast.ToString(user.ID),
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status,
	}
}

// UsersProto convert ent users to proto
func UsersProto(users []*ent.User) []*user_proto.User {
	result := make([]*user_proto.User, len(users))

	for index, user := range users {
		result[index] = UserProto(user)
	}

	return result
}
