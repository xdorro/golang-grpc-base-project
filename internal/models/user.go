package models

import (
	"reflect"

	userpb "github.com/xdorro/base-project-proto/protos/v1/user"
)

// User is a struct that represents a user
type User struct {
	*Common `bson:",inline"`

	Name     string `json:"name,omitempty" bson:"name,omitempty" `
	Email    string `json:"email,omitempty" bson:"email,omitempty" `
	Password string `json:"password" bson:"password,omitempty"`
}

// CollectionName returns the name of the collection from struct name
func (m User) CollectionName() string {
	collName := reflect.TypeOf(m).Name()
	return collName
}

// BeforeCreate is a hook that is called before the creation of the user
func (m *User) BeforeCreate() {
	if m.Common == nil {
		m.Common = &Common{}
	}

	m.Common.BeforeCreate()
}

// UserToProto converts a user to a proto
func (m *User) UserToProto() *userpb.User {
	return &userpb.User{
		Id:    m.ID.Hex(),
		Name:  m.Name,
		Email: m.Email,
	}
}

// UsersToProto converts a list of users to a list of protos
func UsersToProto(users []*User) []*userpb.User {
	result := make([]*userpb.User, len(users))

	for index, user := range users {
		result[index] = user.UserToProto()
	}

	return result
}
