package models

import (
	"reflect"

	userpb "github.com/xdorro/proto-base-project/protos/v1/user"
)

// User is a struct that represents a user
type User struct {
	*Common `bson:",inline"`

	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
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

// ToProto converts a user to a proto
func (m *User) ToProto() *userpb.User {
	return &userpb.User{
		Id:    m.ID.Hex(),
		Name:  m.Name,
		Email: m.Email,
	}
}

// ToUsersProto converts a slice of users to a slice of protos
func ToUsersProto(users []*User) []*userpb.User {
	result := make([]*userpb.User, len(users))

	for index, user := range users {
		result[index] = user.ToProto()
	}

	return result
}
